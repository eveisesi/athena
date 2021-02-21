package cache

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type characterService interface {
	Character(ctx context.Context, id uint) (*athena.Character, error)
	SetCharacter(ctx context.Context, characterID uint, character *athena.Character) error
	Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error)
	SetCharacters(ctx context.Context, records []*athena.Character, operators ...*athena.Operator) error
	CharacterCorporationHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CharacterCorporationHistory, error)
	SetCharacterCorporationHistory(ctx context.Context, records []*athena.CharacterCorporationHistory, operators ...*athena.Operator) error
}

const (
	keyCharacter                   = "athena::character::%d"
	keyCharacters                  = "athena::characters::%x"
	keyCharacterCorporationHistory = "athena::characters::%x::history"
)

func (s *service) Character(ctx context.Context, characterID uint) (*athena.Character, error) {

	key := fmt.Sprintf(keyCharacter, characterID)
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

func (s *service) SetCharacter(ctx context.Context, characterID uint, character *athena.Character) error {

	key := fmt.Sprintf(keyCharacter, characterID)
	data, err := json.Marshal(character)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error) {

	keyData, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCharacters, sha1.Sum(keyData))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errKeyNotFound, key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	records := make([]*athena.Character, 0, len(members))
	for _, member := range members {
		var record = new(athena.Character)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}

		records = append(records, record)

	}

	return records, nil

}

func (s *service) SetCharacters(ctx context.Context, records []*athena.Character, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal characters for cache: %w", err)
		}

		members = append(members, string(data))
	}

	keyData, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCharacters, sha1.Sum(keyData))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache characters %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) CharacterCorporationHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CharacterCorporationHistory, error) {

	keyData, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCharacterCorporationHistory, sha256.Sum224(keyData))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errKeyNotFound, key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	records := make([]*athena.CharacterCorporationHistory, 0, len(members))
	for _, member := range members {
		var record = new(athena.CharacterCorporationHistory)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}

		records = append(records, record)

	}

	return records, nil

}

func (s *service) SetCharacterCorporationHistory(ctx context.Context, records []*athena.CharacterCorporationHistory, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal character corporation history records for cache: %w", err)
		}

		members = append(members, string(data))
	}

	keyData, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCharacterCorporationHistory, sha256.Sum224(keyData))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache corporation history %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}
