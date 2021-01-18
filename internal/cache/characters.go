package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type characterService interface {
	Character(ctx context.Context, id string) (*athena.Character, error)
	SetCharacter(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) error
	Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error)
	SetCharacters(ctx context.Context, operators []*athena.Operator, characters []*athena.Character, optionFuncs ...OptionFunc) error
}

const (
	keyCharacter  = "athena::character::%s"
	keyCharacters = "athena::characters::%s"
)

func (s *service) Character(ctx context.Context, id string) (*athena.Character, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCharacter, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var character = new(athena.Character)

	err = json.Unmarshal([]byte(result), &character)
	if err != nil {
		return nil, err
	}

	return character, nil

}

func (s *service) SetCharacter(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(character)
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCharacter, character.ID.Hex()), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCharacters, fmt.Sprintf("%x", bs))).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var characters = make([]*athena.Character, 0)

	err = json.Unmarshal([]byte(result), &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil

}

func (s *service) SetCharacters(ctx context.Context, operators []*athena.Operator, characters []*athena.Character, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	data, err = json.Marshal(characters)
	if err != nil {
		return fmt.Errorf("Failed to marsahl payload: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCharacters, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
