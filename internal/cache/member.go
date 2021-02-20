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

type memberService interface {
	Member(ctx context.Context, memberID uint) (*athena.Member, error)
	SetMember(ctx context.Context, memberID uint, member *athena.Member) error
	Members(ctx context.Context, operators ...*athena.Operator) ([]*athena.Member, error)
	SetMembers(ctx context.Context, members []*athena.Member, operators ...*athena.Operator) error
}

const (
	keyMember  = "athena::member::%d"
	keyMembers = "athena::members::%x"
)

func (s *service) Member(ctx context.Context, memberID uint) (*athena.Member, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMember, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var member = new(athena.Member)

		err = json.Unmarshal(result, member)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal member onto struct")
		}

		return member, nil
	}

	return nil, nil
}

func (s *service) SetMember(ctx context.Context, memberID uint, member *athena.Member) error {

	data, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMember, memberID), data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) Members(ctx context.Context, operators ...*athena.Operator) ([]*athena.Member, error) {

	keyData, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyMembers, sha1.Sum(keyData))

	values, err := s.client.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var members = make([]*athena.Member, 0)
	for _, v := range values {
		var member = new(athena.Member)
		err = json.Unmarshal([]byte(v), member)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal member from cache: %w", err)
		}

		members = append(members, member)

	}

	return members, nil

}

// SetMembers caches a slice of members using the slice of operators used to fetch that slice of members.
func (s *service) SetMembers(ctx context.Context, members []*athena.Member, operators ...*athena.Operator) error {

	values := make([]string, 0, len(members))
	for _, member := range members {
		data, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("Failed to marsahl member for cache: %w", err)
		}

		values = append(values, string(data))
	}

	keyData, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyMembers, sha1.Sum(keyData))
	_, err = s.client.SAdd(ctx, key, values).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	_, err = s.client.Expire(ctx, key, time.Minute*20).Result()
	if err != nil {
		return fmt.Errorf("failed to set expiry for key %s: %w", key, err)
	}

	return nil
}
