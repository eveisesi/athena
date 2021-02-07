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
	MySQL struct {
		Host string `required:"true"`
		Port int    `required:"true"`
		User string `required:"true"`
		Pass string `required:"true"`
		DB   string `required:"true"`
	}

	Redis struct {
		Host string `required:"true"`
		Port uint   `required:"true"`
	}

	NewRelic struct {
		AppName string `envconfig:"NEW_RELIC_APP_NAME"`
	}

	Env athena.Environment `default:"development"`

	Developer struct {
		Name string
	}

	Log struct {
		Level string `required:"true"`
	}

	Server struct {
		Port uint `required:"true"`
	}

	Auth struct {
		ClientID         string `required:"true"`
		ClientSecret     string `required:"true"`
		RedirectURL      string `required:"true"`
		AuthorizationURL string `required:"true"`
		TokenURL         string `required:"true"`
		JWKSURL          string `required:"true"`
	}

	UserAgent string `required:"true"`
}

func (c config) validateEnvironment() bool {
	return c.Env.Validate()
}
