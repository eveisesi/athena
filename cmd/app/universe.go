package main

import (
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/urfave/cli"
)

func universeCommand(c *cli.Context) error {

	basics := basics("universe")

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	universe := universe.NewService(basics.logger, cache, esi, basics.repositories.universe)

	err = universe.InitializeUniverse()
	if err != nil {
		basics.logger.WithError(err).Error("failed to initialize the universe")
		return err
	}

	basics.logger.Info("universe initialize successfully")

	return nil

}
