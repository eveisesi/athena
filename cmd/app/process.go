package main

import (
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/mongodb"
	"github.com/eveisesi/athena/internal/processor"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/urfave/cli"
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

	cloneRepo, err := mongodb.NewCloneRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize clone repository")
	}

	universeRepo, err := mongodb.NewUniverseRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize universe repository")
	}

	cache := cache.NewService(basics.redis)
	esi := esi.NewService(cache, basics.client, basics.cfg.UserAgent)

	character := character.NewService(cache, esi, characterRepo)
	corporation := corporation.NewService(cache, esi, corporationRepo)
	alliance := alliance.NewService(cache, esi, allianceRepo)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	universe := universe.NewService(basics.logger, cache, esi, universeRepo)
	member := member.NewService(auth, cache, alliance, character, corporation, memberRepo)
	location := location.NewService(basics.logger, cache, esi, universe, locationRepo)
	clone := clone.NewService(basics.logger, cache, esi, universe, cloneRepo)

	processor := processor.NewService(basics.logger, cache, esi, member, location, clone)

	processor.Run()

	return nil

}
