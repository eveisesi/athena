package main

import (
	"fmt"
	"log"
	"math/rand"
	nethttp "net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	mpb "github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
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
		{
			Name:   "server",
			Usage:  "Initialize the HTTP Services responsible for handling HTTP Requests to this application",
			Action: serverCommand,
		},
		{
			Name:   "processor",
			Action: processorCommand,
		},
		{
			Name:   "universe",
			Action: universeCommand,
		},
		{
			Name: "progress",
			Action: func(c *cli.Context) error {
				var wg sync.WaitGroup
				// pass &wg (optional), so p will wait for it eventually
				p := mpb.New(mpb.WithWaitGroup(&wg))
				total, numBars := 100, 3
				wg.Add(numBars)

				for i := 0; i < numBars; i++ {
					name := fmt.Sprintf("Bar#%d:", i)
					bar := p.AddBar(int64(total),
						mpb.PrependDecorators(
							// simple name decorator
							decor.Name(name),
							// decor.DSyncWidth bit enables column width synchronization
							decor.Percentage(decor.WCSyncSpace),
						),
						mpb.AppendDecorators(
							// replace ETA decorator with "done" message, OnComplete event
							decor.OnComplete(
								// ETA decorator with ewma age of 60
								decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
							),
						),
					)
					// simulating some work
					go func() {
						defer wg.Done()
						rng := rand.New(rand.NewSource(time.Now().UnixNano()))
						max := 100 * time.Millisecond
						for i := 0; i < total; i++ {
							// start variable is solely for EWMA calculation
							// EWMA's unit of measure is an iteration's duration
							start := time.Now()
							time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
							bar.Increment()
							// we need to call DecoratorEwmaUpdate to fulfill ewma decorator's contract
							bar.DecoratorEwmaUpdate(time.Since(start))
						}
					}()
				}
				// Waiting for passed &wg and for all bars to complete and flush
				p.Wait()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
