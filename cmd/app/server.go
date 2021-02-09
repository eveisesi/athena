package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/member"

	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"

	"github.com/eveisesi/athena/internal/server"
	"github.com/urfave/cli"
)

func serverCommand(c *cli.Context) error {

	basics := basics("server")

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	// universe := universe.NewService(basics.logger, cache, esi, basics.repositories.universe)
	// location := location.NewService(basics.logger, cache, esi, universe, basics.repositories.location)
	// clone := clone.NewService(basics.logger, cache, esi, universe, basics.repositories.clone)
	corporation := corporation.NewService(cache, esi, basics.repositories.corporation)
	alliance := alliance.NewService(cache, esi, basics.repositories.alliance)
	character := character.NewService(basics.logger, cache, esi, corporation, basics.repositories.character)
	// contact := contact.NewService(basics.logger, cache, esi, etag, universe, basics.repositories.contact)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	member := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)

	server := server.NewServer(
		basics.cfg.Server.Port,
		basics.cfg.Env,
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
