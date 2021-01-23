package main

import (
	"github.com/eveisesi/athena"
	"golang.org/x/oauth2"
)

func getAuthConfig(cfg config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Auth.ClientID,
		ClientSecret: cfg.Auth.ClientSecret,
		RedirectURL:  cfg.Auth.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Auth.AuthorizationURL,
			TokenURL: cfg.Auth.TokenURL,
		},
		Scopes: []string{
			athena.ReadLocationV1,
		},
	}
}
