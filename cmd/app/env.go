package main

import (
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func loadConfig() (cfg config, err error) {
	_ = godotenv.Load(".config/app.env")

	err = envconfig.Process("", &cfg)
	if err != nil {
		return config{}, err
	}

	if !cfg.validateEnvironment() {
		return config{}, fmt.Errorf("invalid env %s declared", cfg.Env)
	}

	return

}

type config struct {
	Mongo struct {
		Host     string
		Port     int
		User     string
		Pass     string
		Name     string
		AuthMech string `default:"SCRAM-SHA-256"`
	}

	Redis struct {
		Host string
		Port uint
	}

	NewRelic struct {
		Enabled bool   `envconfig:"NEW_RELIC_ENABLED" default:"false"`
		AppName string `envconfig:"NEW_RELIC_APP_NAME"`
	}

	Env athena.Environment `default:"production"`

	Developer struct {
		Name string
	}

	Log struct {
		Level string
	}

	Server struct {
		Port uint
	}

	Auth struct {
		ClientID         string
		ClientSecret     string
		RedirectURL      string
		AuthorizationURL string
		TokenURL         string
		JWKSURL          string
	}
}

func (c config) validateEnvironment() bool {
	return c.Env.Validate()
}
