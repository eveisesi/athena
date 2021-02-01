package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type characterService interface {
	Character(ctx context.Context, id uint) (*athena.Character, error)
	SetCharacter(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) error
	// Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error)
	// SetCharacters(ctx context.Context, operators []*athena.Operator, characters []*athena.Character, optionFuncs ...OptionFunc) error
}

const (
	keyCharacter = "athena::character::%d"
	// keyCharacters = "athena::characters::%d"
)

func (s *service) Character(ctx context.Context, id uint) (*athena.Character, error) {

	key := fmt.Sprintf(keyCharacter, id)
	result, err := s.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var character = new(athena.Character)
	err = json.Unmarshal([]byte(result), &character)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return character, nil

}

func (s *service) SetCharacter(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyCharacter, character.ID)
	data, err := json.Marshal(character)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

// func (s *service) Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error) {

// 	data, err := json.Marshal(operators)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal operators: %w", err)
// 	}

// 	h := sha256.New()
// 	_, _ = h.Write(data)
// 	bs := h.Sum(nil)

// 	result, err := s.client.Get(ctx, fmt.Sprintf(keyCharacters, fmt.Sprintf("%x", bs))).Result()
// 	if err != nil && err != redis.Nil {
// 		return nil, err
// 	}

// 	if len(result) == 0 {
// 		return nil, nil
// 	}

// 	var characters = make([]*athena.Character, 0)

// 	err = json.Unmarshal([]byte(result), &characters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return characters, nil

// }

// func (s *service) SetCharacters(ctx context.Context, operators []*athena.Operator, characters []*athena.Character, optionFuncs ...OptionFunc) error {

// 	data, err := json.Marshal(operators)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal operators: %w", err)
// 	}

// 	h := sha256.New()
// 	_, _ = h.Write(data)
// 	bs := h.Sum(nil)

// 	data, err = json.Marshal(characters)
// 	if err != nil {
// 		return fmt.Errorf("Failed to marsahl payload: %w", err)
// 	}

// 	options := applyOptionFuncs(nil, optionFuncs)

// 	_, err = s.client.Set(ctx, fmt.Sprintf(keyCharacters, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
// 	if err != nil {
// 		return fmt.Errorf("failed to write to cache: %w", err)
// 	}

// 	return nil

// }
