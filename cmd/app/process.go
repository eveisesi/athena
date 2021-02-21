package main

import (
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/asset"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/contact"
	"github.com/eveisesi/athena/internal/contract"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/mail"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/processor"
	"github.com/eveisesi/athena/internal/skill"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/eveisesi/athena/internal/wallet"
	"github.com/urfave/cli"
)

func processorCommand(c *cli.Context) error {

	basics := basics("processor")

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	alliance := alliance.NewService(basics.logger, cache, esi, basics.repositories.alliance)
	corporation := corporation.NewService(basics.logger, cache, esi, alliance, basics.repositories.corporation)
	character := character.NewService(basics.logger, cache, esi, corporation, basics.repositories.character)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	universe := universe.NewService(basics.logger, cache, esi, basics.repositories.universe)

	asset := asset.NewService(basics.logger, cache, esi, universe, basics.repositories.asset)
	member := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)
	clone := clone.NewService(basics.logger, cache, esi, universe, basics.repositories.clone)
	contact := contact.NewService(basics.logger, cache, esi, universe, alliance, character, corporation, basics.repositories.contact)
	contract := contract.NewService(basics.logger, cache, esi, universe, alliance, character, corporation, basics.repositories.contract)
	location := location.NewService(basics.logger, cache, esi, universe, basics.repositories.location)
	mail := mail.NewService(basics.logger, cache, esi, character, alliance, corporation, basics.repositories.mail)
	skill := skill.NewService(basics.logger, cache, esi, etag, universe, basics.repositories.skill)
	wallet := wallet.NewService(basics.logger, cache, esi, universe, alliance, corporation, character, basics.repositories.wallet)

	processor := processor.NewService(basics.logger, cache, member)

	processor.SetScopeMap(
		buildScopeMap(
			location, clone, contact,
			mail, skill, wallet,
			asset, contract,
		),
	)

	processor.Run()

	return nil

}
