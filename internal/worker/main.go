package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	hapi "github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	hworker "github.com/Donders-Institute/filer-gateway/internal/worker/handler"
	log "github.com/sirupsen/logrus"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	redisAddr   *string
	nworkers    *int
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	nworkers = flag.Int("p", 4, "`number` of concurrent workers per queue")
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
	fmt.Printf("\nBackground task worker for filer gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {
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
	}, bokchoy.WithMaxRetries(3), bokchoy.WithRetryIntervals([]time.Duration{
		30 * time.Second,
		60 * time.Second,
		120 * time.Second,
	}), bokchoy.WithLogger(logger), bokchoy.WithTTL(7*24*time.Hour))

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	bok.Use(middleware.Recoverer)
	bok.Use(middleware.DefaultLogger)

	// add handler to handle tasks in the queue of `hapi.QueueSetProject`
	bok.Queue(hapi.QueueSetProject).Handle(
		&hworker.SetProjectResourceHandler{},
		bokchoy.WithConcurrency(*nworkers),
	)

	// add handler to handle tasks in the queue of `hapi.QueueSetUser`
	bok.Queue(hapi.QueueSetUser).Handle(
		&hworker.SetUserResourceHandler{},
		bokchoy.WithConcurrency(*nworkers),
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Print("Received signal, gracefully stopping")
			bok.Stop(ctx)
		}
	}()

	bok.Run(ctx)
}
