package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type locationService interface {
	MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error)
	SetMemberLocation(ctx context.Context, memberID uint, location *athena.MemberLocation, optionFuncs ...OptionFunc) error
	MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error)
	SetMemberOnline(ctx context.Context, memberID uint, online *athena.MemberOnline, optionFuncs ...OptionFunc) error
	MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error)
	SetMemberShip(ctx context.Context, memberID uint, ship *athena.MemberShip, optionFuncs ...OptionFunc) error
}

const (
	keyMemberLocation = "athena::member::%d::location"
	keyMemberOnline   = "athena::member::%d::online"
	keyMemberShip     = "athena::member::%d::ship"
)

func (s *service) MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error) {

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

func (s *service) SetMemberLocation(ctx context.Context, memberID uint, location *athena.MemberLocation, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberLocation, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error) {

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

func (s *service) SetMemberOnline(ctx context.Context, memberID uint, online *athena.MemberOnline, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(online)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberOnline, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error) {

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

func (s *service) SetMemberShip(ctx context.Context, memberID uint, ship *athena.MemberShip, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(ship)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberShip, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
