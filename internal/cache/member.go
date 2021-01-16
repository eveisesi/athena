package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type memberRepository struct {
	client *redis.Client
	expiry time.Duration
}

const MEMBER = "athena::member::%s"
const MEMBERID = "athena::member::id::%d"
const MEMBERS = "athena::members::%s"

func NewMemberRepository(client *redis.Client, expiry time.Duration) athena.CacheMemberRepository {

	return &memberRepository{
		client,
		expiry,
	}
}

func (r *memberRepository) Members(ctx context.Context, operators []*athena.Operator) ([]*athena.Member, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := r.client.Get(ctx, fmt.Sprintf(MEMBERS, fmt.Sprintf("%x", bs))).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) > 0 {
		var members = make([]*athena.Member, 0)

		err = json.Unmarshal([]byte(result), &members)
		if err != nil {
			return nil, err
		}

		return members, nil
	}

	return nil, nil

}

func (r *memberRepository) SetMembers(ctx context.Context, operators []*athena.Operator, members []*athena.Member) error {

	opts, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(opts)
	bs := h.Sum(nil)

	key := fmt.Sprintf(MEMBERS, fmt.Sprintf("%x", bs))

	data, err := json.Marshal(members)
	if err != nil {
		return fmt.Errorf("Failed to marsahl payload: %w", err)
	}

	_, err = r.client.Set(ctx, key, data, r.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}
