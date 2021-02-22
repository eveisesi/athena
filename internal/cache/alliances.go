package cache

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type allianceService interface {
	Alliance(ctx context.Context, id uint) (*athena.Alliance, error)
	SetAlliance(ctx context.Context, allianceID uint, alliance *athena.Alliance) error
	Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error)
	SetAlliances(ctx context.Context, alliances []*athena.Alliance, operators ...*athena.Operator) error
}

const (
	keyAlliance  = "athena::alliance::%d"
	keyAlliances = "athena::alliances::%x"
)

func (s *service) Alliance(ctx context.Context, allianceID uint) (*athena.Alliance, error) {

	key := fmt.Sprintf(keyAlliance, allianceID)
	result, err := s.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
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

func (s *service) SetAlliance(ctx context.Context, allianceID uint, alliance *athena.Alliance) error {

	key := fmt.Sprintf(keyAlliance, allianceID)

	data, err := json.Marshal(alliance)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyAlliances, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var alliances = make([]*athena.Alliance, 0, len(members))
	for _, member := range members {
		var alliance = new(athena.Alliance)
		err = json.Unmarshal([]byte(member), alliance)
		if err != nil {
			return nil, err
		}

		alliances = append(alliances, alliance)
	}

	return alliances, nil

}

func (s *service) SetAlliances(ctx context.Context, alliances []*athena.Alliance, operators ...*athena.Operator) error {

	members := make([]string, 0, len(alliances))
	for _, alliance := range alliances {
		data, err := json.Marshal(alliance)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal alliance: %w", err)
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

	key := fmt.Sprintf(keyAlliances, sha1.Sum(keyData))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	_, err = s.client.Expire(ctx, key, time.Minute*10).Result()
	if err != nil {
		return fmt.Errorf(errFailedToSetExpiry, key, err)
	}

	return nil

}
