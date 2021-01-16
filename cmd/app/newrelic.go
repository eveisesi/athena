package main

import (
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func loadNewrelicApplication(cfg config, logger *logrus.Logger) (app *newrelic.Application, err error) {

	appName := cfg.NewRelic.AppName

	if cfg.Env != athena.ProductionEnvironment {
		appName = fmt.Sprintf("%s-%s", cfg.Env, appName)
		if cfg.Developer.Name != "" {
			appName = fmt.Sprintf("%s-%s", cfg.Developer.Name, appName)
		}
	}

	opts := []newrelic.ConfigOption{}
	opts = append(opts, newrelic.ConfigFromEnvironment())
	opts = append(opts, newrelic.ConfigAppName(appName))
	app, err = newrelic.NewApplication(opts...)
	if err != nil {
		return nil, err
	}

	err = app.WaitForConnection(time.Second * 20)

	return

}
