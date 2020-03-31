package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/config"
	"github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"

	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	optsPort    *int
	redisAddr   *string
	configFile  *string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	redisAddr = flag.String("r", "redis:6379", "redis service `address`")
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

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	api := operations.NewFilerGatewayAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalf("%s", err)
		}
	}()

	server.Port = *optsPort

	// initiate blochy queue for setting project roles
	var logger logging.Logger

	ctx := context.Background()
	bok, err := bokchoy.New(ctx, bokchoy.Config{
		Broker: bokchoy.BrokerConfig{
			Type: "redis",
			Redis: bokchoy.RedisConfig{
				Type: "client",
				Client: bokchoy.RedisClientConfig{
					Addr: *redisAddr,
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

	// associate handler functions with implementations
	api.GetTasksTypeIDHandler = operations.GetTasksTypeIDHandlerFunc(handler.GetTask(ctx, bok))
	api.PostProjectsHandler = operations.PostProjectsHandlerFunc(handler.CreateProject(ctx, bok))
	api.PatchProjectsIDHandler = operations.PatchProjectsIDHandlerFunc(handler.UpdateProject(ctx, bok))
	api.GetProjectsIDHandler = operations.GetProjectsIDHandlerFunc(handler.GetProjectResource())
	api.GetProjectsIDMembersHandler = operations.GetProjectsIDMembersHandlerFunc(handler.GetProjectMembers())
	api.GetProjectsIDStorageHandler = operations.GetProjectsIDStorageHandlerFunc(handler.GetProjectStorage())
	api.GetUsersIDHandler = operations.GetUsersIDHandlerFunc(handler.GetUserResource())

	// configure API
	server.ConfigureAPI()

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalf("%s", err)
	}
}
