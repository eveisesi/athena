package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
	"github.com/sirkon/go-format"
)

type characterService interface {
	Character(ctx context.Context, id uint) (*athena.Character, error)
	SetCharacter(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) error
	CharacterCorporationHistory(ctx context.Context, id uint) ([]*athena.CharacterCorporationHistory, error)
	SetCharacterCorporationHistory(ctx context.Context, id uint, history []*athena.CharacterCorporationHistory, optionFuncs ...OptionFunc) error
	// Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error)
	// SetCharacters(ctx context.Context, operators []*athena.Operator, characters []*athena.Character, optionFuncs ...OptionFunc) error
}

const (
	keyCharacter                   = "athena::character::${id}"
	keyCharacterCorporationHistory = "athena::character::${id}::history"
	// keyCharacters = "athena::characters::%d"
)

func (s *service) Character(ctx context.Context, id uint) (*athena.Character, error) {

	key := format.Formatm(keyCharacter, format.Values{
		"id": id,
	})
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

	key := format.Formatm(keyCharacter, format.Values{
		"id": character.ID,
	})
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

func (s *service) CharacterCorporationHistory(ctx context.Context, id uint) ([]*athena.CharacterCorporationHistory, error) {

	key := format.Formatm(keyCharacterCorporationHistory, format.Values{
		"id": id,
	})
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errKeyNotFound, key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	records := make([]*athena.CharacterCorporationHistory, len(members))
	for i, member := range members {
		var record = new(athena.CharacterCorporationHistory)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}

		records[i] = record

	}

	return records, nil

}

func (s *service) SetCharacterCorporationHistory(ctx context.Context, id uint, records []*athena.CharacterCorporationHistory, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]string, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal character corporation history records for cache: %w", err)
		}

		members = append(members, string(data))
	}

	key := format.Formatm(keyCharacterCorporationHistory, format.Values{
		"id": id,
	})
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache corporation history for character %d: %w", id, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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
