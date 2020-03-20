package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/loads"

	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	log "github.com/sirupsen/logrus"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	optsPort    *int
	redisAddr   *string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	redisAddr = flag.String("r", "redis:6379", "redis service `address`")

	flag.Usage = usage

	flag.Parse()

	// set logging
	log.SetOutput(os.Stderr)

	// set logging level
	llevel := log.InfoLevel
	if *optsVerbose {
		llevel = log.DebugLevel
	}
	log.SetLevel(llevel)
}

func usage() {
	fmt.Printf("\nAPI server of filer gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewFilerGatewayAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
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
		log.Fatalln(err)
	}
}
