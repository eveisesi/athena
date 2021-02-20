package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type memberService interface {
	Member(ctx context.Context, memberID uint) (*athena.Member, error)
	SetMember(ctx context.Context, memberID uint, member *athena.Member, optionFuncs ...OptionFunc) error
	Members(ctx context.Context, operators []*athena.Operator) ([]*athena.Member, error)
	SetMembers(ctx context.Context, operators []*athena.Operator, members []*athena.Member, optionFuncs ...OptionFunc) error
}

const (
	keyMember  = "athena::member::%d"
	keyMembers = "athena::members::%s"
)

func (s *service) Member(ctx context.Context, memberID uint) (*athena.Member, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMember, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var member = new(athena.Member)

		err = json.Unmarshal(result, member)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal member onto struct")
		}

		return member, nil
	}

	return nil, nil
}

func (s *service) SetMember(ctx context.Context, memberID uint, member *athena.Member, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMember, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) Members(ctx context.Context, operators []*athena.Operator) ([]*athena.Member, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMembers, fmt.Sprintf("%x", bs))).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var members = make([]*athena.Member, 0)

		err = json.Unmarshal([]byte(result), &members)
		if err != nil {
			return nil, err
		}

		return members, nil
	}

	return nil, nil

}

// SetMembers caches a slice of members using the slice of operators used to fetch that slice of members.
func (s *service) SetMembers(ctx context.Context, operators []*athena.Operator, members []*athena.Member, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	data, err = json.Marshal(members)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMembers, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}
