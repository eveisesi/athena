package member

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/lestrrat-go/jwx/jwt"
)

type Service interface {
	Member(ctx context.Context, memberID uint) (*athena.Member, error)
	UpdateMember(ctx context.Context, member *athena.Member) (*athena.Member, error)
	Login(ctx context.Context, code, state string) error
	ValidateToken(ctx context.Context, member *athena.Member) (*athena.Member, error)
	Middleware(next http.Handler) http.Handler
	MemberFromToken(ctx context.Context, token jwt.Token) (*athena.Member, error)
	MemberFromContext(ctx context.Context) *athena.Member
}

type service struct {
	auth        auth.Service
	cache       cache.Service
	character   character.Service
	corporation corporation.Service
	alliance    alliance.Service

	member athena.MemberRepository
}

type ctxKey struct {
	name string
}

// const (
// 	serviceIdentifier = "Member Service"
// )

var userCtxKey = ctxKey{name: "user"}

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

// func (s *service) ExpiredTokens(ctx context.Context) ([]*athena.Member, error) {

// 	members, err := s.member.Members(ctx, athena.NewLessThanOperator("expires", time.Now()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	spew.Dump(members)

// }

func (s *service) ValidateToken(ctx context.Context, member *athena.Member) (*athena.Member, error) {

	currentToken := member.AccessToken
	member, err := s.auth.ValidateToken(ctx, member)
	if err != nil {
		return nil, err
	}

	if member.AccessToken != currentToken {
		_, err = s.member.UpdateMember(ctx, member.ID, member)
		if err != nil {
			return nil, err
		}
	}

	return member, nil

}

func (s *service) Member(ctx context.Context, memberID uint) (*athena.Member, error) {

	member, err := s.cache.Member(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if member != nil {
		return member, nil
	}

	member, err = s.member.Member(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return member, nil

}

func (s *service) UpdateMember(ctx context.Context, member *athena.Member) (*athena.Member, error) {
	return s.member.UpdateMember(ctx, member.ID, member)
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
		return fmt.Errorf("failed to exchange state and code for token: %w", err)
	}

	token, err := s.auth.ParseAndVerifyToken(ctx, bearer.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to parse and/or verify token: %w", err)
	}

	member, err := s.MemberFromToken(ctx, token)
	if err != nil {
		return err
	}

	member.AccessToken.SetValid(bearer.AccessToken)
	member.RefreshToken.SetValid(bearer.RefreshToken)
	member.Expires.SetValid(bearer.Expiry)

	_, err = s.member.UpdateMember(ctx, member.ID, member)
	if err != nil {
		return err
	}

	s.cache.PushIDToProcessorQueue(ctx, member.ID)
	_ = s.cache.SetMember(ctx, member.ID, member)

	attempt.Status = athena.CompletedAuthStatus
	attempt.Token = member.AccessToken

	_, err = s.auth.UpdateAuthAttempt(ctx, attempt.State, attempt)
	if err != nil {
		return err
	}

	return nil

}

func (s *service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		token := r.Header.Get("authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		token = token[7:]
		parsed, err := s.auth.ParseAndVerifyToken(ctx, token)
		if err != nil {
			fmt.Printf("Error Parsing Token: %v, stop auth, proceeding with request\n", err)
			next.ServeHTTP(w, r)
			return
		}

		sub := parsed.Subject()
		if sub == "" {
			fmt.Printf("Invalid Token Subject: %s\n", sub)
			next.ServeHTTP(w, r)
			return
		}

		memberID, err := memberIDFromSubject(sub)
		if err != nil {
			fmt.Printf("unable to parse character ID from token sub: %s: %s\n", sub, err)
			next.ServeHTTP(w, r)
			return
		}

		member, err := s.member.Member(ctx, memberID)
		if err != nil {
			fmt.Printf("failed to retrieve member with ID from token: %s\n", err)
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, member)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (s *service) MemberFromContext(ctx context.Context) *athena.Member {

	member := ctx.Value(userCtxKey)
	if member == nil {
		return nil
	}

	if _, ok := member.(*athena.Member); !ok {
		return nil
	}

	return member.(*athena.Member)

}

func (s *service) MemberFromToken(ctx context.Context, token jwt.Token) (*athena.Member, error) {

	sub := token.Subject()
	if sub == "" {
		return nil, fmt.Errorf("unexpected empty subject")
	}

	memberID, err := memberIDFromSubject(sub)
	if err != nil {
		return nil, err
	}

	member, err := s.member.Member(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if member == nil {
		// This is a new member, lets create a record for them.
		_, err := s.character.FetchCharacter(ctx, memberID)
		if err != nil {
			return nil, err
		}

		character, err := s.character.Character(ctx, memberID)
		if err != nil {
			return nil, err
		}

		_, err = s.corporation.FetchCorporation(ctx, character.CorporationID)
		if err != nil {
			return nil, err
		}

		corporation, err := s.corporation.Corporation(ctx, character.CorporationID)
		if err != nil {
			return nil, err
		}

		if corporation.AllianceID.Valid {
			_, err = s.alliance.FetchAlliance(ctx, corporation.AllianceID.Uint)
			if err != nil {
				return nil, err
			}
		}

		member = &athena.Member{
			ID:        character.ID,
			LastLogin: time.Now(),
			IsNew:     true,
		}
	}

	claims := token.PrivateClaims()

	if _, ok := claims["owner"]; !ok {
		return nil, fmt.Errorf("failed to process token. owner hash is missing")
	}

	member.OwnerHash.SetValid(claims["owner"].(string))

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

	if member.IsNew {
		member, err = s.member.CreateMember(ctx, member)
		if err != nil {
			return nil, err
		}
	} else {
		member, err = s.member.UpdateMember(ctx, member.ID, member)
		if err != nil {
			return nil, err
		}
	}

	return member, nil

}

func memberIDFromSubject(sub string) (uint, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	id, err := strconv.ParseUint(parts[2], 10, 32)

	return uint(id), err

}
