package main

import (
	"database/sql"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/mysqldb"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	logger *logrus.Logger
	err    error
)

type app struct {
	cfg          config
	newrelic     *newrelic.Application
	logger       *logrus.Logger
	db           *sql.DB
	redis        *redis.Client
	client       *nethttp.Client
	repositories repositories
}

type repositories struct {
	alliance    athena.AllianceRepository
	character   athena.CharacterRepository
	contact     athena.MemberContactRepository
	corporation athena.CorporationRepository
	clone       athena.CloneRepository
	etag        athena.EtagRepository
	location    athena.MemberLocationRepository
	member      athena.MemberRepository
	migration   athena.MigrationRepository
	skill       athena.MemberSkillRepository
	universe    athena.UniverseRepository
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

	app.db, err = makeMySQL(app.cfg)
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

	app.repositories = repositories{
		member:      mysqldb.NewMemberRepository(app.db),
		character:   mysqldb.NewCharacterRepository(app.db),
		corporation: mysqldb.NewCorporationRepository(app.db),
		alliance:    mysqldb.NewAllianceRepository(app.db),
		etag:        mysqldb.NewEtagRepository(app.db),
		universe:    mysqldb.NewUniverseRepository(app.db),
		migration:   mysqldb.NewMigrationRepository(app.db),
		clone:       mysqldb.NewCloneRepository(app.db),
		skill:       mysqldb.NewSkillRepository(app.db),
		// location:    location,
		// contact:     contact,
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
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "skip",
					Usage: "skip certain imports so we don't have to wait to load the entire universe",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "disable the progress bar so debug output doesn't get overwritten",
				},
			},
		},
		{
			Name: "migrate",
			Subcommands: []cli.Command{
				{
					Name:   "up",
					Action: migrateUpCommand,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "steps",
							Value: -1,
						},
					},
				},
				{
					Name:   "down",
					Action: migrateDownCommand,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "steps",
							Value: -1,
						},
					},
				},
				{
					Name:   "create",
					Action: migrateCreateCommand,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "name",
							Required: true,
						},
					},
				},
			},
		},
		{
			Name:  "token",
			Usage: "Commands for managing tokens",
			Subcommands: []cli.Command{
				{
					Name:   "manual",
					Usage:  "Provide a MemberID and that ID will be pushed to the Processor Queue",
					Action: manuallyPushIDToQueue,
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:     "id",
							Required: true,
						},
					},
				},
				{
					Name:   "add",
					Usage:  "Will parse and create an account for the prompted access token and refresh token",
					Action: addMemberByCLI,
				},
				{
					Name:   "refresh",
					Action: refreshMemberToken,
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:     "id",
							Required: true,
						},
					},
				},
			},
		},
		{
			Name:   "test",
			Action: testCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
