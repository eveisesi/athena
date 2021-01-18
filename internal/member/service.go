package member

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/volatiletech/null"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/lestrrat-go/jwx/jwt"
)

type Service interface {
	Login(ctx context.Context, code, state string) error
}

type service struct {
	auth        auth.Service
	cache       cache.Service
	character   character.Service
	corporation corporation.Service
	alliance    alliance.Service

	member athena.MemberRepository
}

func NewService(auth auth.Service, cache cache.Service, alliance alliance.Service, character character.Service, corporation corporation.Service, member athena.MemberRepository) Service {
	return &service{
		auth:        auth,
		cache:       cache,
		character:   character,
		corporation: corporation,
		alliance:    alliance,

		member: member,
	}
}

func (s *service) Login(ctx context.Context, code, state string) error {

	attempt, err := s.auth.AuthAttempt(ctx, state)
	if err != nil {
		return err
	}
	if attempt != nil && attempt.Status == athena.InvalidAuthStatus {
		return fmt.Errorf("attempt is no longer valid")
	}

	bearer, err := s.auth.BearerForCode(ctx, code)
	if err != nil {
		msg := "failed to exchange state and code for token"
		return fmt.Errorf("%s: %w", msg, err)
	}

	token, err := s.auth.ParseAndVerifyToken(ctx, bearer.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to parse and/or verify token: %w", err)
	}

	member, err := s.memberFromToken(ctx, token)
	if err != nil {
		return err
	}

	member.AccessToken = bearer.AccessToken
	member.RefreshToken = bearer.RefreshToken
	member.Expires = bearer.Expiry

	if member.IsNew {
		s.cache.PushIDToProcessorQueue(ctx, member.ID)
	}

	_, err = s.member.UpdateMember(ctx, member.ID.Hex(), member)
	if err != nil {
		return err
	}

	attempt.Status = athena.CompletedAuthStatus
	attempt.Token = null.NewString(member.AccessToken, true)

	_, err = s.auth.UpdateAuthAttempt(ctx, attempt.State, attempt)
	if err != nil {
		return err
	}

	return nil

}

func (s *service) memberFromToken(ctx context.Context, token jwt.Token) (*athena.Member, error) {

	sub := token.Subject()
	if sub == "" {
		return nil, fmt.Errorf("unexpected empty subject")
	}

	memberID, err := memberIDFromSubject(sub)
	if err != nil {
		return nil, err
	}

	operators := []*athena.Operator{athena.NewEqualOperator("character_id", memberID), athena.NewLimitOperator(1)}

	members, err := s.cache.Members(ctx, operators)
	if err != nil {
		return nil, err
	}

	if len(members) > 1 {
		return nil, fmt.Errorf("invalid number of members returned")
	}

	if members != nil {
		return members[0], nil
	}

	members, err = s.member.Members(ctx, operators...)
	if err != nil {
		return nil, err
	}

	if len(members) == 1 {

		_ = s.cache.SetMembers(ctx, operators, members)

		return members[0], nil
	}

	// This is a new member, lets create a record for them.
	character, err := s.character.CharacterByCharacterID(ctx, memberID, character.NewOptionFuncs())
	if err != nil {
		return nil, err
	}

	corporation, err := s.corporation.CorporationByCorporationID(ctx, character.CorporationID, corporation.NewOptionFuncs())
	if err != nil {
		return nil, err
	}

	if corporation.AllianceID.Valid {
		_, err = s.alliance.AllianceByAllianceID(ctx, corporation.AllianceID.Uint, alliance.NewOptionFuncs())
		if err != nil {
			return nil, err
		}
	}

	claims := token.PrivateClaims()

	member := &athena.Member{
		CharacterID: character.CharacterID,
		LastLogin:   time.Now(),
		IsNew:       true,
	}

	if _, ok := claims["owner"]; !ok {
		return nil, fmt.Errorf("failed to process token. owner hash is missing")
	}

	member.OwnerHash = claims["owner"].(string)

	if _, ok := claims["scp"]; ok {
		scp := []athena.MemberScope{}
		switch a := claims["scp"].(type) {
		case []interface{}:
			for _, v := range a {
				scp = append(scp, athena.MemberScope{
					Scope: athena.Scope(v.(string)),
				})
			}
		case string:
			scp = append(scp, athena.MemberScope{
				Scope: athena.Scope(a),
			})
		}

		member.Scopes = scp
	}

	member, err = s.member.CreateMember(ctx, member)
	if err != nil {
		return nil, err
	}

	s.cache.PushIDToProcessorQueue(ctx, member.ID)

	return member, nil

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}
