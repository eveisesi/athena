package main

import (
	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/mongodb"
	"github.com/eveisesi/athena/internal/processor"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func processorCommand(c *cli.Context) error {

	basics := basics("processor")

	locationRepo, err := mongodb.NewLocationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatalln("failed to initialize location repository")
	}

	memberRepo, err := mongodb.NewMemberRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize member repository")
	}

	characterRepo, err := mongodb.NewCharacterRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize character repository")
	}

	corporationRepo, err := mongodb.NewCorporationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize corporation repository")
	}

	allianceRepo, err := mongodb.NewAllianceRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize alliance repository")
	}

	cache := cache.NewService(basics.redis)
	esi := esi.NewService(cache, basics.client, basics.cfg.UserAgent)

	corporation := corporation.NewService(cache, esi, corporationRepo)
	alliance := alliance.NewService(cache, esi, allianceRepo)

	character := character.NewService(cache, esi, characterRepo)

	auth := auth.NewService(
		cache,
		&oauth2.Config{
			ClientID:     basics.cfg.Auth.ClientID,
			ClientSecret: basics.cfg.Auth.ClientSecret,
			RedirectURL:  basics.cfg.Auth.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  basics.cfg.Auth.AuthorizationURL,
				TokenURL: basics.cfg.Auth.TokenURL,
			},
			Scopes: []string{athena.READ_SHIP_V1},
		},
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	member := member.NewService(auth, cache, alliance, character, corporation, memberRepo)
	location := location.NewService(basics.logger, cache, esi, locationRepo)

	processor := processor.NewService(basics.logger, cache, esi, member, location)

	processor.Run()

	return nil

}
