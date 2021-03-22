package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	hapi "github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	hworker "github.com/Donders-Institute/filer-gateway/internal/worker/handler"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	redisURL    *string
	nworkers    *int
	configFile  *string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	nworkers = flag.Int("p", 4, "`number` of global concurrent workers")
	redisURL = flag.String("r", "redis://redis:6379", "redis service `address`")
	configFile = flag.String("c", os.Getenv("FILER_GATEWAY_WORKER_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/worker.log",
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
	fmt.Printf("\nTask worker of filer gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {
	// initiate blochy queue for setting project roles
	var logger logging.Logger

	// redis client instance for notifying cache update
	redisOpts, err := redis.ParseURL(*redisURL)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx := context.Background()

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
	}, bokchoy.WithLogger(logger), bokchoy.WithTTL(7*24*time.Hour), bokchoy.WithConcurrency(*nworkers))

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	bok.Use(middleware.Recoverer)
	bok.Use(middleware.DefaultLogger)

	// add handler to handle tasks in the queue of `hapi.QueueSetProject`
	bok.Queue(hapi.QueueSetProject).Handle(
		&hworker.SetProjectResourceHandler{
			ConfigFile:        *configFile,
			ApiNotifierClient: redis.NewClient(redisOpts),
		},
	)

	// add handler to handle tasks in the queue of `hapi.QueueSetUser`
	bok.Queue(hapi.QueueSetUser).Handle(
		&hworker.SetUserResourceHandler{
			ConfigFile:        *configFile,
			ApiNotifierClient: redis.NewClient(redisOpts),
		},
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Warnf("Received signal, gracefully stopping")
			bok.Stop(ctx)
		}
	}()

	bok.Run(ctx)
}
