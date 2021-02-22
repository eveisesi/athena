package main

import (
	"context"
	"fmt"
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
	"github.com/volatiletech/null"
	"golang.org/x/oauth2"
)

func manuallyPushIDToQueue(c *cli.Context) error {

	basics := basics("token-refresh")
	var ctx = context.Background()

	cache := cache.NewService(basics.redis)

	memberID := c.Int64("id")
	members, err := basics.repositories.member.Members(ctx, athena.NewEqualOperator("id", memberID))
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

	memberServ := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)

	memberID := c.Int64("id")
	members, err := basics.repositories.member.Members(ctx, athena.NewEqualOperator("id", memberID))
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

	if c.Bool("reset") {
		for i := range member.Scopes {
			member.Scopes[i].Expiry = null.TimeFromPtr(nil)
		}
		_, err := basics.repositories.member.UpdateMember(ctx, member.ID, member)
		if err != nil {
			basics.logger.WithError(err).Fatal("failed to update member")
		}
	}

	fmt.Println("Member Token validated successfully")

	return nil

}

func resetMemberByCLI(c *cli.Context) error {

	basics := basics("token-reset")
	var ctx = context.Background()

	memberID := c.Int64("id")
	members, err := basics.repositories.member.Members(ctx, athena.NewEqualOperator("id", memberID))
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to fetch member from db")
	}

	if len(members) != 1 {
		basics.logger.WithField("count", len(members)).Fatal("unexpected number of member returned")
	}

	member := members[0]

	for i := range member.Scopes {
		member.Scopes[i].Expiry = null.TimeFromPtr(nil)
	}
	_, err = basics.repositories.member.UpdateMember(ctx, member.ID, member)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to update member")
	}

	fmt.Println("Member Scopes reset successfully")

	return nil

}

func addMemberByCLI(c *cli.Context) error {
	basics := basics("token-refresh")

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

	memberServ := member.NewService(auth, cache, alliance, character, corporation, basics.repositories.member)

	ctx := context.Background()

	// reader := bufio.NewReader(os.Stdin)

	// fmt.Print("Access Token: ")
	// at, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print("Refresh Token: ")
	// rt, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// at = strings.TrimSuffix(at, "\n")
	// rt = strings.TrimSuffix(rt, "\n")

	oauth2Token := new(oauth2.Token)
	oauth2Token.AccessToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkpXVC1TaWduYXR1cmUtS2V5IiwidHlwIjoiSldUIn0.eyJzY3AiOlsiZXNpLWFzc2V0cy5yZWFkX2Fzc2V0cy52MSIsImVzaS1jaGFyYWN0ZXJzLnJlYWRfY29udGFjdHMudjEiLCJlc2ktY2xvbmVzLnJlYWRfY2xvbmVzLnYxIiwiZXNpLWNsb25lcy5yZWFkX2ltcGxhbnRzLnYxIiwiZXNpLWNvbnRyYWN0cy5yZWFkX2NoYXJhY3Rlcl9jb250cmFjdHMudjEiLCJlc2ktbG9jYXRpb24ucmVhZF9sb2NhdGlvbi52MSIsImVzaS1sb2NhdGlvbi5yZWFkX3NoaXBfdHlwZS52MSIsImVzaS1tYWlsLnJlYWRfbWFpbC52MSIsImVzaS1za2lsbHMucmVhZF9za2lsbHF1ZXVlLnYxIiwiZXNpLXNraWxscy5yZWFkX3NraWxscy52MSIsImVzaS11bml2ZXJzZS5yZWFkX3N0cnVjdHVyZXMudjEiLCJlc2ktd2FsbGV0LnJlYWRfY2hhcmFjdGVyX3dhbGxldC52MSJdLCJqdGkiOiI5MjkyYTM3Mi0zNzE0LTQ2MjMtOGVkNy0zNDNmZTUyYmU0YzkiLCJraWQiOiJKV1QtU2lnbmF0dXJlLUtleSIsInN1YiI6IkNIQVJBQ1RFUjpFVkU6OTAyMTg5MjMiLCJhenAiOiI2MWMyMTJiNjFjMWM0NDVmYTUzZmNjYWI2YmJkZDZiMyIsInRlbmFudCI6InRyYW5xdWlsaXR5IiwidGllciI6ImxpdmUiLCJyZWdpb24iOiJ3b3JsZCIsIm5hbWUiOiIxMjNuaWNrIiwib3duZXIiOiJDVzJzU09wN2x3MEVBTlVCcTd5QkVqcFdkWkk9IiwiZXhwIjoxNjEzODY2MjMwLCJpc3MiOiJsb2dpbi5ldmVvbmxpbmUuY29tIn0.aCRi7NB2TY4D5L0i_wHbtVkx9LAj17dmVq9MIeNBZdSGAz_t-e7tDXnRRpdHEcCOMIQT5e4AKu9G35jgCHqpEmjUvv1IwYaybMhlaykZhlViVOoU9r-ysDeKDIepv4AzrXn50MOuEAEQZg6lG3SEq8naHEF7xFfGqIEjJKffe9_GMm3sVau6wZ08sYE1rbFUTO2_qR4zGiGnsb5GVxrrnYOZVR2yDuzw2TfzAM75jyNMMv3Lr0-opjdsj2RwAcp87Z9fURtj_XwenUDb-pBF_gQR-FG0EbNcQ9YpTTN8PDUI14aZ5j3Dp-A__12EThW58nmLuPQS6pm0AM-EkC8XtQ"
	oauth2Token.RefreshToken = "hb0K7o32K0y82Jo4EQ8r4g=="
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

	member.AccessToken.SetValid(oauth2Token.AccessToken)
	member.Expires.SetValid(oauth2Token.Expiry)
	member.RefreshToken.SetValid(oauth2Token.RefreshToken)

	_, err = basics.repositories.member.UpdateMember(ctx, member.ID, member)
	if err != nil {
		return err
	}

	skipQueue := c.Bool("skipQueue")

	if !skipQueue {
		fmt.Println("Skip Queue False, push to queue")
		cache.PushIDToProcessorQueue(ctx, member.ID)
	}

	_ = cache.SetMember(ctx, member.ID, member)

	return nil

}
