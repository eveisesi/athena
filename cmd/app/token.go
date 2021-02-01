package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/member"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func manuallyPushIDToQueue(c *cli.Context) error {

	basics := basics("token-refresh")
	var ctx = context.Background()

	cache := cache.NewService(basics.redis)

	memberID := c.Int64("id")
	members, err := basics.repositories.member.Members(ctx, athena.NewEqualOperator("character_id", memberID))
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to fetch member from db")
	}

	if len(members) != 1 {
		basics.logger.WithField("count", len(members)).Fatal("unexpected number of member returned")
	}

	member := members[0]

	cache.PushIDToProcessorQueue(ctx, member.ID)

	return nil

}

func refreshMemberToken(c *cli.Context) error {

	basics := basics("token-refresh")
	var ctx = context.Background()

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	character := character.NewService(cache, esi, basics.repositories.character)
	corporation := corporation.NewService(cache, esi, basics.repositories.corporation)
	alliance := alliance.NewService(cache, esi, basics.repositories.alliance)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	memberServ := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)

	memberID := c.Int64("id")
	members, err := basics.repositories.member.Members(ctx, athena.NewEqualOperator("character_id", memberID))
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to fetch member from db")
	}

	if len(members) != 1 {
		basics.logger.WithField("count", len(members)).Fatal("unexpected number of member returned")
	}

	member := members[0]

	_, err = memberServ.ValidateToken(ctx, member)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to validate token")
	}

	fmt.Println("Member Token validated successfully")

	return nil

}

func addMemberByCLI(c *cli.Context) error {
	basics := basics("token-refresh")

	cache := cache.NewService(basics.redis)
	etag := etag.NewService(basics.logger, cache, basics.repositories.etag)
	esi := esi.NewService(basics.client, cache, etag, basics.cfg.UserAgent)

	character := character.NewService(cache, esi, basics.repositories.character)
	corporation := corporation.NewService(cache, esi, basics.repositories.corporation)
	alliance := alliance.NewService(cache, esi, basics.repositories.alliance)

	auth := auth.NewService(
		cache,
		getAuthConfig(basics.cfg),
		basics.client,
		basics.cfg.Auth.JWKSURL,
	)

	memberServ := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)

	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Access Token: ")
	at, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Refresh Token: ")
	rt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	at = strings.TrimSuffix(at, "\n")
	rt = strings.TrimSuffix(rt, "\n")

	oauth2Token := new(oauth2.Token)
	oauth2Token.AccessToken = at
	oauth2Token.RefreshToken = rt
	oauth2Token.Expiry = time.Now()

	oauth := getAuthConfig(basics.cfg)

	tokenSource := oauth.TokenSource(ctx, oauth2Token)
	newOauth2Token, err := tokenSource.Token()
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to refresh token")
	}

	token, err := auth.ParseAndVerifyToken(ctx, newOauth2Token.AccessToken)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to parse token")
	}

	member, err := memberServ.MemberFromToken(ctx, token)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to create member from token")
	}

	member.AccessToken = oauth2Token.AccessToken
	member.RefreshToken = oauth2Token.RefreshToken
	member.Expires = oauth2Token.Expiry

	_, err = basics.repositories.member.UpdateMember(ctx, member.ID, member)
	if err != nil {
		return err
	}

	cache.PushIDToProcessorQueue(ctx, member.ID)
	_ = cache.SetMember(ctx, member.ID, member)

	return nil

}