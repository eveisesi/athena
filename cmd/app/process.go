package main

import (
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/contact"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/processor"
	"github.com/eveisesi/athena/internal/skill"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/urfave/cli"
)

func processorCommand(c *cli.Context) error {

	basics := basics("processor")

	cache := cache.NewService(basics.redis)
	esi := esi.NewService(cache, basics.client, basics.cfg.UserAgent)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)

	character := character.NewService(cache, esi, basics.repositories.character)
	corporation := corporation.NewService(cache, esi, basics.repositories.corporation)
	alliance := alliance.NewService(cache, esi, basics.repositories.alliance)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	universe := universe.NewService(basics.logger, cache, esi, basics.repositories.universe)
	member := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)
	location := location.NewService(basics.logger, cache, esi, universe, basics.repositories.location)
	clone := clone.NewService(basics.logger, cache, esi, universe, basics.repositories.clone)
	contact := contact.NewService(basics.logger, cache, esi, etag, universe, alliance, character, corporation, basics.repositories.contact)
	skill := skill.NewService(basics.logger, cache, esi, etag, universe, basics.repositories.skill)

	processor := processor.NewService(basics.logger, cache, member)

	processor.SetScopeMap(buildScopeMap(location, clone, contact, skill))

	listScopes()

	processor.Run()

	return nil

}
