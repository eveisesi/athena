package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/athena"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/mongodb"

	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"golang.org/x/oauth2"

	"github.com/eveisesi/athena/internal/server"
	"github.com/urfave/cli"
)

func serverCommand(c *cli.Context) error {

	basics := basics("server")

	cache := cache.NewService(basics.redis)

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
			Scopes: []string{athena.ReadShipV1.String()},
		},
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	member := member.NewService(auth, cache, alliance, character, corporation, memberRepo)

	server := server.NewServer(
		basics.cfg.Server.Port,
		basics.cfg.Env,
		basics.db,
		basics.logger,
		cache,
		basics.newrelic,
		auth,
		member,
		character,
		corporation,
		alliance,
	)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Run()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		basics.logger.WithError(err).Fatal("server encountered an unexpected error and had to quit")
	case sig := <-osSignals:
		basics.logger.WithField("sig", sig).Info("interrupt signal received, starting server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = server.GracefullyShutdown(ctx)
		if err != nil {
			basics.logger.WithError(err).Fatal("failed to shutdown server")
		}

		basics.logger.Info("server gracefully shutdown successfully")
	}

	return nil

}
