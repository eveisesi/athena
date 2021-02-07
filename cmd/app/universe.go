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

	skip := c.StringSlice("skip")
	debug := c.Bool("debug")

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	universeServ := universe.NewService(basics.logger, cache, esi, basics.repositories.universe)

	opts := make([]universe.OptionFunc, 0)
	if len(skip) > 0 {
		if contains(skip, "chr") {
			opts = append(opts, universe.WithoutChr())
		}
		if contains(skip, "loc") {
			opts = append(opts, universe.WithoutLoc())
		}
		if contains(skip, "inv") {
			opts = append(opts, universe.WithoutInv())
		}
	}

	if debug {
		opts = append(opts, universe.WithDisableProgress())
	}

	err = universeServ.InitializeUniverse(opts...)
	if err != nil {
		basics.logger.WithError(err).Error("failed to initialize the universe")
		return err
	}

	basics.logger.Info("universe initialize successfully")

	return nil

}

func contains(slc []string, s string) bool {
	for _, v := range slc {
		if s == v {
			return true
		}
	}

	return false
}
