package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type locationService interface {
	MemberLocation(ctx context.Context, memberID string) (*athena.MemberLocation, error)
	SetMemberLocation(ctx context.Context, memberID string, location *athena.MemberLocation, optionFuncs ...OptionFunc) error
	MemberOnline(ctx context.Context, memberID string) (*athena.MemberOnline, error)
	SetMemberOnline(ctx context.Context, memberID string, online *athena.MemberOnline, optionFuncs ...OptionFunc) error
	MemberShip(ctx context.Context, memberID string) (*athena.MemberShip, error)
	SetMemberShip(ctx context.Context, memberID string, ship *athena.MemberShip, optionFuncs ...OptionFunc) error
}

const (
	keyMemberLocation = "athena::member::%s::location"
	keyMemberOnline   = "athena::member::%s::online"
	keyMemberShip     = "athena::member::%s::ship"
)

func (s *service) MemberLocation(ctx context.Context, memberID string) (*athena.MemberLocation, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberLocation, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var location = new(athena.MemberLocation)

	err = json.Unmarshal(result, location)
	if err != nil {
		return nil, err
	}

	return location, nil

}

func (s *service) SetMemberLocation(ctx context.Context, memberID string, location *athena.MemberLocation, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberLocation, fmt.Sprintf("%s", memberID)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberOnline(ctx context.Context, memberID string) (*athena.MemberOnline, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberOnline, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var online = new(athena.MemberOnline)

	err = json.Unmarshal(result, online)
	if err != nil {
		return nil, err
	}

	return online, nil

}

func (s *service) SetMemberOnline(ctx context.Context, memberID string, online *athena.MemberOnline, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(online)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberOnline, fmt.Sprintf("%s", memberID)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberShip(ctx context.Context, memberID string) (*athena.MemberShip, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberShip, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var ship = new(athena.MemberShip)

	err = json.Unmarshal(result, ship)
	if err != nil {
		return nil, err
	}

	return ship, nil

}

func (s *service) SetMemberShip(ctx context.Context, memberID string, ship *athena.MemberShip, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(ship)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberShip, fmt.Sprintf("%s", memberID)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
