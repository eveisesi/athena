package cache

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type corporationService interface {
	Corporation(ctx context.Context, corporationID uint) (*athena.Corporation, error)
	SetCorporation(ctx context.Context, corporationID uint, corporation *athena.Corporation) error

	Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error)
	SetCorporations(ctx context.Context, corporations []*athena.Corporation, operators ...*athena.Operator) error

	CorporationAllianceHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CorporationAllianceHistory, error)
	SetCorporationAllianceHistory(ctx context.Context, records []*athena.CorporationAllianceHistory, operators ...*athena.Operator) error
}

const (
	keyCorporation                = "athena::corporation::%d"
	keyCorporationAllianceHistory = "athena::corporation::%x::history"
	keyCorporations               = "athena::corporations::%x"
)

func (s *service) Corporation(ctx context.Context, corporationID uint) (*athena.Corporation, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCorporation, corporationID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var corporation = new(athena.Corporation)
	err = json.Unmarshal([]byte(result), &corporation)
	if err != nil {
		return nil, err
	}

	return corporation, nil

}

func (s *service) SetCorporation(ctx context.Context, corporationID uint, corporation *athena.Corporation) error {

	data, err := json.Marshal(corporation)
	if err != nil {
		return fmt.Errorf("failed to marshal corporation: %w", err)
	}

	key := fmt.Sprintf(keyCorporation, corporationID)
	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error) {

	keyData, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyCorporations, sha1.Sum(keyData))

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(members) == 0 || errors.Is(err, redis.Nil) {
		return nil, nil
	}

	records := make([]*athena.Corporation, 0, len(members))
	for _, member := range members {
		var record = new(athena.Corporation)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}
		records = append(records, record)

	}

	return records, nil

}

func (s *service) SetCorporations(ctx context.Context, corporations []*athena.Corporation, operators ...*athena.Operator) error {

	members := make([]string, 0, len(corporations))
	for _, corporation := range corporations {
		data, err := json.Marshal(corporation)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal corporation: %w", err)
		}

		members = append(members, string(data))
	}

	if len(operators) == 0 {
		return fmt.Errorf("length of operators should be greater 0")
	}

	keyData, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyCorporations, sha1.Sum(keyData))

	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) CorporationAllianceHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CorporationAllianceHistory, error) {

	keyData, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCorporationAllianceHistory, sha1.Sum(keyData))

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errKeyNotFound, key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	records := make([]*athena.CorporationAllianceHistory, 0, len(members))
	for _, member := range members {
		var record = new(athena.CorporationAllianceHistory)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}

		records = append(records, record)

	}

	return records, nil

}

func (s *service) SetCorporationAllianceHistory(ctx context.Context, records []*athena.CorporationAllianceHistory, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	keyData, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators for key: %w", err)
	}

	key := fmt.Sprintf(keyCorporationAllianceHistory, sha1.Sum(keyData))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache alliance history for corporations %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to set expiry on key %s: %w", key, err)
	}

	return nil

}
