package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dccn-tg/filer-gateway/internal/api-server/config"
	"github.com/dccn-tg/filer-gateway/internal/api-server/handler"
	"github.com/dccn-tg/filer-gateway/pkg/swagger/server/models"
	"github.com/dccn-tg/filer-gateway/pkg/swagger/server/restapi"
	"github.com/dccn-tg/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/s12v/go-jwks"
	"github.com/square/go-jose"

	"github.com/hurngchunlee/bokchoy"
	"github.com/hurngchunlee/bokchoy/logging"
	"github.com/hurngchunlee/bokchoy/middleware"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	optsPort    *int
	redisURL    *string
	configFile  *string

	pcache handler.ProjectResourceCache
	ucache handler.UserResourceCache
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	redisURL = flag.String("r", "redis://redis:6379", "redis service `address`")
	configFile = flag.String("c", os.Getenv("FILER_GATEWAY_APISERVER_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/api-server.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

func usage() {
	fmt.Printf("\nAPI server of filer gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	// load global configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", *configFile)
	}

	// redis client instance for notifying cache update
	redisOpts, err := redis.ParseURL(*redisURL)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// initialize Cache
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// project cache
	ppubsub := redis.NewClient(redisOpts).Subscribe(ctx, "api_pcache_update")
	pcache = handler.ProjectResourceCache{
		Config:   cfg,
		Context:  ctx,
		Notifier: ppubsub.Channel(),
	}
	pcache.Init()

	// user cache
	upubsub := redis.NewClient(redisOpts).Subscribe(ctx, "api_ucache_update")
	ucache = handler.UserResourceCache{
		Config:   cfg,
		Context:  ctx,
		Notifier: upubsub.Channel(),
	}
	ucache.Init()

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	api := operations.NewFilerGatewayAPI(swaggerSpec)
	api.UseRedoc()
	server := restapi.NewServer(api)

	// actions to take when the main program exists.
	defer func() {
		// stop the redis Pub/Sub for cache refresh notification.
		ppubsub.Close()

		// stop all background services of the context.
		cancel()

		// stop API server.
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalf("%s", err)
		}
	}()

	server.Port = *optsPort
	server.ListenLimit = 10
	server.TLSListenLimit = 10

	// initiate blochy queue for setting project roles
	var logger logging.Logger

	bok, err := bokchoy.New(ctx, bokchoy.Config{
		Broker: bokchoy.BrokerConfig{
			Type: "redis",
			Redis: bokchoy.RedisConfig{
				Type: "client",
				Client: bokchoy.RedisClientConfig{
					Addr: redisOpts.Addr,
				},
			},
		},
	}, bokchoy.WithLogger(logger), bokchoy.WithTTL(7*24*time.Hour))

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	bok.Use(middleware.Recoverer)
	bok.Use(middleware.DefaultLogger)

	// authentication with api key.
	api.APIKeyHeaderAuth = func(token string) (*models.Principle, error) {

		if token != cfg.ApiKey {
			return nil, errors.New(401, "incorrect api key auth")
		}

		// there is no user information attached, set the principle as empty string.
		principle := models.Principle("")
		return &principle, nil
	}

	// authentication with username/password.
	api.BasicAuthAuth = func(username, password string) (*models.Principle, error) {

		pass, ok := cfg.Auth[username]

		if !ok || pass != password {
			return nil, errors.New(401, "incorrect username/password")
		}

		// there is login user information attached, set the principle as the username.
		principle := models.Principle(username)
		return &principle, nil
	}

	// authentication with oauth2 token.
	api.Oauth2Auth = func(tokenStr string, scopes []string) (*models.Principle, error) {

		// custom claims data structure, this should match the
		// data structure expected from the authentication server.
		type IDServerClaims struct {
			Scope    []string `json:"scope"`
			Audience []string `json:"aud"`
			ClientID string   `json:"client_id"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(tokenStr, &IDServerClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New(401, "unexpected signing method: %v", token.Header["alg"])
			}

			// get public key from the auth server
			// TODO: discover jwks endpoint using oidc client.
			jwksSource := jwks.NewWebSource(cfg.JwksEndpoint)
			jwksClient := jwks.NewDefaultClient(
				jwksSource,
				time.Hour,    // Refresh keys every 1 hour
				12*time.Hour, // Expire keys after 12 hours
			)

			var jwk *jose.JSONWebKey
			jwk, err := jwksClient.GetEncryptionKey(token.Header["kid"].(string))
			if err != nil {
				return nil, errors.New(401, "cannot retrieve encryption key: %s", err)
			}

			return jwk.Key, nil
		})

		if err != nil {
			return nil, errors.New(401, "invalid token: %s", err)
		}

		// check token scope
		claims, ok := token.Claims.(*IDServerClaims)
		if !ok {
			return nil, errors.New(401, "cannot get claims from the token")
		}

		inScope := func(target string) bool {
			for _, s := range claims.Scope {
				if s == target {
					return true
				}
			}
			return false
		}

		for _, scope := range scopes {
			if !inScope(scope) {
				return nil, errors.New(401, "token not in scope: %s", scope)
			}
		}

		principle := models.Principle(claims.ClientID)
		return &principle, nil
	}

	// associate handler functions with implementations
	api.GetPingHandler = operations.GetPingHandlerFunc(handler.GetPing(cfg))

	api.GetMetricsHandler = operations.GetMetricsHandlerFunc(handler.GetMetrics(&ucache, &pcache))

	api.GetTasksTypeIDHandler = operations.GetTasksTypeIDHandlerFunc(handler.GetTask(ctx, bok))

	api.GetProjectsIDHandler = operations.GetProjectsIDHandlerFunc(handler.GetProjectResource(&pcache))
	// api.GetProjectsIDMembersHandler = operations.GetProjectsIDMembersHandlerFunc(handler.GetProjectMembers())
	// api.GetProjectsIDStorageHandler = operations.GetProjectsIDStorageHandlerFunc(handler.GetProjectStorage())

	api.GetProjectsHandler = operations.GetProjectsHandlerFunc(handler.GetProjects(&pcache))
	api.PostProjectsHandler = operations.PostProjectsHandlerFunc(handler.CreateProject(ctx, bok))
	api.PatchProjectsIDHandler = operations.PatchProjectsIDHandlerFunc(handler.UpdateProject(ctx, bok))

	api.GetUsersHandler = operations.GetUsersHandlerFunc(handler.GetUsers(&ucache, &pcache))
	api.GetUsersIDHandler = operations.GetUsersIDHandlerFunc(handler.GetUserResource(&ucache, &pcache))
	api.PostUsersHandler = operations.PostUsersHandlerFunc(handler.CreateUserResource(ctx, bok))
	api.PatchUsersIDHandler = operations.PatchUsersIDHandlerFunc(handler.UpdateUserResource(ctx, bok))
	api.DeleteUsersIDHandler = operations.DeleteUsersIDHandlerFunc(handler.DeleteUserResource(ctx, bok))

	// configure API
	server.ConfigureAPI()

	// Start API server
	if err := server.Serve(); err != nil {
		log.Fatalf("%s", err)
	}
}
