package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type cloneService interface {
	MemberClones(ctx context.Context, memberID string) (*athena.MemberClones, error)
	SetMemberClones(ctx context.Context, memberID string, clones *athena.MemberClones, optionFuncs ...OptionFunc) error
	MemberImplants(ctx context.Context, memberID string) (*athena.MemberImplants, error)
	SetMemberImplants(ctx context.Context, memberID string, implants *athena.MemberImplants, optionFuncs ...OptionFunc) error
}

const (
	keyMemberClones   = "athena::member::%s::clones"
	keyMemberImplants = "athena::member::%s::implants"
)

func (s *service) MemberClones(ctx context.Context, memberID string) (*athena.MemberClones, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberClones, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var clones = new(athena.MemberClones)

	err = json.Unmarshal(result, clones)
	if err != nil {
		return nil, err
	}

	return clones, nil

}

func (s *service) SetMemberClones(ctx context.Context, memberID string, clones *athena.MemberClones, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(clones)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberClones, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberImplants(ctx context.Context, memberID string) (*athena.MemberImplants, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberImplants, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var implants = new(athena.MemberImplants)

	err = json.Unmarshal(result, &implants)
	if err != nil {
		return nil, err
	}

	return implants, nil

}

func (s *service) SetMemberImplants(ctx context.Context, memberID string, implants *athena.MemberImplants, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(implants)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberImplants, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
