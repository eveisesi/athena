package main

import (
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	logger *logrus.Logger
	err    error
)

type app struct {
	cfg      config
	newrelic *newrelic.Application
	logger   *logrus.Logger
	db       *mongo.Database
	redis    *redis.Client
	client   *nethttp.Client
}

// basics initializes the following
// loadConfig - parses environment variables and applies them to a struct
// loadLogger - takes in a configuration and intializes a logrus logger
// loadDB - takes in a configuration and establishes a connection with our datastore, in this application that is mongoDB
// loadRedis - takes in a configuration and establises a connection with our cache, in this application that is Redis
// loadNewrelic - takes in a configuration and configures a NR App to report metrics to NewRelic for monitoring
// loadClient - create a client from the net/http library that is used on all outgoing http requests
func basics(command string) *app {

	app := app{}

	app.cfg, err = loadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	app.logger, err = loadLogger(app.cfg, command)
	if err != nil {
		log.Fatalf("failed to load logger: %s", err)
	}

	app.cfg.NewRelic.AppName = fmt.Sprintf("%s-%s", app.cfg.NewRelic.AppName, command)

	app.newrelic, err = loadNewrelicApplication(app.cfg, app.logger)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure NR App")
	}

	app.db, err = makeMongoDB(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to make mongo db connection")
	}

	app.redis = makeRedis(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure redis client")
	}

	app.client = &nethttp.Client{
		Timeout:   time.Second * 5,
		Transport: newrelic.NewRoundTripper(nil),
	}
	return &app

}

func main() {
	app := cli.NewApp()
	app.Name = "Athena CLI"
	app.UsageText = "athena"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "server",
			Usage:  "Initialize the HTTP Services responsible for handling HTTP Requests to this application",
			Action: serverCommand,
		},
		cli.Command{
			Name:   "processor",
			Action: processorCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
