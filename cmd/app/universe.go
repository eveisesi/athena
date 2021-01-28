package main

import (
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/mongodb"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/urfave/cli"
)

func universeCommand(c *cli.Context) error {

	basics := basics("universe")

	universeRepo, err := mongodb.NewUniverseRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize universe repository")
	}

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	universe := universe.NewService(basics.logger, cache, esi, universeRepo)

	err = universe.InitializeUniverse()
	if err != nil {
		basics.logger.WithError(err).Error("failed to initialize the universe")
		return err
	}

	basics.logger.Info("universe initialize successfully")

	return nil

}
