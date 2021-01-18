package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type allianceService interface {
	Alliance(ctx context.Context, id string) (*athena.Alliance, error)
	SetAlliance(ctx context.Context, alliance *athena.Alliance, optionFuncs ...OptionFunc) error
	Alliances(ctx context.Context, operators []*athena.Operator) ([]*athena.Alliance, error)
	SetAlliances(ctx context.Context, operators []*athena.Operator, alliances []*athena.Alliance, optionFuncs ...OptionFunc) error
}

const (
	keyAlliance  = "athena::alliance::%s"
	keyAlliances = "athena::alliances::%s"
)

func (s *service) Alliance(ctx context.Context, id string) (*athena.Alliance, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyAlliance, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var alliance = new(athena.Alliance)

	err = json.Unmarshal([]byte(result), &alliance)
	if err != nil {
		return nil, err
	}

	return alliance, nil

}

func (s *service) SetAlliance(ctx context.Context, alliance *athena.Alliance, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(alliance)
	if err != nil {
		return fmt.Errorf("failed to marshal alliance: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyAlliance, alliance.ID.Hex()), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Alliances(ctx context.Context, operators []*athena.Operator) ([]*athena.Alliance, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := s.client.Get(ctx, fmt.Sprintf(keyAlliances, fmt.Sprintf("%x", bs))).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var alliances = make([]*athena.Alliance, 0)

		err = json.Unmarshal([]byte(result), &alliances)
		if err != nil {
			return nil, err
		}

		return alliances, nil
	}

	return nil, nil

}

func (s *service) SetAlliances(ctx context.Context, operators []*athena.Operator, alliances []*athena.Alliance, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	data, err = json.Marshal(alliances)
	if err != nil {
		return fmt.Errorf("Failed to marsahl payload: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyAlliances, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
