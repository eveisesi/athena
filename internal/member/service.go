package member

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/lestrrat-go/jwx/jwt"
)

type Service interface {
	Login(ctx context.Context, code, state string) error
	// ProcessUserByToken(token jwt.Token)
}

type service struct {
	auth      auth.Service
	cache     cache.Service
	character character.Service

	member athena.MemberRepository
}

func NewService(auth auth.Service, cache cache.Service, character character.Service, member athena.MemberRepository) Service {
	return &service{
		auth:      auth,
		cache:     cache,
		character: character,

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

	spew.Dump(bearer, token, member, err)

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

	// character, err := s.character.Character(ctx, memberID)

	return nil, nil

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}
