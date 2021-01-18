package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type corporationService interface {
	Corporation(ctx context.Context, id string) (*athena.Corporation, error)
	SetCorporation(ctx context.Context, corporation *athena.Corporation, optionFuncs ...OptionFunc) error
	Corporations(ctx context.Context, operators []*athena.Operator) ([]*athena.Corporation, error)
	SetCorporations(ctx context.Context, operators []*athena.Operator, corporations []*athena.Corporation, optionFuncs ...OptionFunc) error
}

const (
	keyCorporation  = "athena::corporation::%s"
	keyCorporations = "athena::corporations::%s"
)

func (s *service) Corporation(ctx context.Context, id string) (*athena.Corporation, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCorporation, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var corporation = new(athena.Corporation)

	err = json.Unmarshal([]byte(result), &corporation)
	if err != nil {
		return nil, err
	}

	return corporation, nil

}

func (s *service) SetCorporation(ctx context.Context, corporation *athena.Corporation, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(corporation)
	if err != nil {
		return fmt.Errorf("failed to marshal corporation: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCorporation, corporation.ID.Hex()), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Corporations(ctx context.Context, operators []*athena.Operator) ([]*athena.Corporation, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCorporations, fmt.Sprintf("%x", bs))).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var corporations = make([]*athena.Corporation, 0)

		err = json.Unmarshal([]byte(result), &corporations)
		if err != nil {
			return nil, err
		}

		return corporations, nil
	}

	return nil, nil

}

func (s *service) SetCorporations(ctx context.Context, operators []*athena.Operator, corporations []*athena.Corporation, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	data, err = json.Marshal(corporations)
	if err != nil {
		return fmt.Errorf("Failed to marsahl payload: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCorporations, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
