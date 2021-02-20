package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirkon/go-format"
)

type cloneService interface {
	MemberClones(ctx context.Context, memberID uint) (*athena.MemberClones, error)
	SetMemberClones(ctx context.Context, memberID uint, clones *athena.MemberClones) error
	MemberImplants(ctx context.Context, memberID uint) ([]*athena.MemberImplant, error)
	SetMemberImplants(ctx context.Context, memberID uint, implants []*athena.MemberImplant) error
}

const (
	keyMemberClone    = "athena::member::${memberID}::clone"
	keyMemberImplants = "athena::member::${memberID}::implants"
)

func (s *service) MemberClones(ctx context.Context, memberID uint) (*athena.MemberClones, error) {

	key := format.Formatm(keyMemberClone, format.Values{
		"memberID": memberID,
	})
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var clone = new(athena.MemberClones)
	err = json.Unmarshal(result, clone)
	if err != nil {
		return nil, err
	}

	return clone, nil

}

func (s *service) SetMemberClones(ctx context.Context, memberID uint, clone *athena.MemberClones) error {

	data, err := json.Marshal(clone)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	key := format.Formatm(keyMemberClone, format.Values{
		"memberID": memberID,
	})
	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberImplants(ctx context.Context, memberID uint) ([]*athena.MemberImplant, error) {

	key := format.Formatm(keyMemberImplants, format.Values{
		"memberID": memberID,
	})
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	implants := make([]*athena.MemberImplant, len(members))
	for i, member := range members {
		var implant = new(athena.MemberImplant)
		err = json.Unmarshal([]byte(member), implant)
		if err != nil {
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s on struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		implants[i] = implant

	}

	return implants, nil

}

func (s *service) SetMemberImplants(ctx context.Context, memberID uint, implants []*athena.MemberImplant) error {

	members := make([]string, 0, len(implants))
	for _, implant := range implants {
		data, err := json.Marshal(implant)
		if err != nil {
			return fmt.Errorf("failed to marshal implants for cache: %w", err)
		}

		members = append(members, string(data))
	}

	key := format.Formatm(keyMemberImplants, format.Values{
		"memberID": memberID,
	})
	_, err := s.client.SAdd(ctx, key, members, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache: %w", err)
	}

	return nil

}
