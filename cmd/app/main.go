package main

import (
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"time"

	"github.com/eveisesi/athena"
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
	cfg          config
	newrelic     *newrelic.Application
	logger       *logrus.Logger
	db           *mongo.Database
	redis        *redis.Client
	client       *nethttp.Client
	repositories repositories
}

type repositories struct {
	member      athena.MemberRepository
	character   athena.CharacterRepository
	corporation athena.CorporationRepository
	alliance    athena.AllianceRepository
	clone       athena.CloneRepository
	etag        athena.EtagRepository
	location    athena.MemberLocationRepository
	universe    athena.UniverseRepository
	contact     athena.MemberContactRepository
	skill       athena.MemberSkillRepository
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

	// app.db, err = makeMongoDB(app.cfg)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to make mongo db connection")
	// }

	app.redis = makeRedis(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure redis client")
	}

	app.client = &nethttp.Client{
		Timeout:   time.Second * 5,
		Transport: newrelic.NewRoundTripper(nil),
	}

	// member, err := mongodb.NewMemberRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize member repository")
	// }

	// character, err := mongodb.NewCharacterRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize character repository")
	// }

	// corporation, err := mongodb.NewCorporationRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize corporation repository")
	// }

	// alliance, err := mongodb.NewAllianceRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize alliance repository")
	// }

	// clone, err := mongodb.NewCloneRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize alliance repository")
	// }

	// etag, err := mongodb.NewEtagRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize etag repository")
	// }

	// location, err := mongodb.NewLocationRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize location repository")
	// }

	// universe, err := mongodb.NewUniverseRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize universe repository")
	// }

	// contact, err := mongodb.NewMemberContactRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize contact repository")
	// }

	// skill, err := mongodb.NewMemberSkillRepository(app.db)
	// if err != nil {
	// 	app.logger.WithError(err).Fatal("failed to initialize skill repository")
	// }

	// app.repositories = repositories{
	// 	member:      member,
	// 	character:   character,
	// 	corporation: corporation,
	// 	alliance:    alliance,
	// 	clone:       clone,
	// 	etag:        etag,
	// 	location:    location,
	// 	universe:    universe,
	// 	contact:     contact,
	// 	skill:       skill,
	// }

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
		// {
		// 	Name:   "test",
		// 	Action: testCommand,
		// },
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
