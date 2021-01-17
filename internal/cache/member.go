package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type memberService interface {
	Members(ctx context.Context, operators []*athena.Operator) ([]*athena.Member, error)
	SetMembers(ctx context.Context, operators []*athena.Operator, members []*athena.Member, optionFuncs ...OptionFunc) error
}

const MEMBER = "athena::member::%s"
const MEMBERID = "athena::member::id::%d"
const MEMBERS = "athena::members::%s"

func (s *service) Members(ctx context.Context, operators []*athena.Operator) ([]*athena.Member, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	result, err := s.client.Get(ctx, fmt.Sprintf(MEMBERS, fmt.Sprintf("%x", bs))).Result()
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

// SetMembers caches a slice of members using the slice of operators used to fetch that slice of members.
func (s *service) SetMembers(ctx context.Context, operators []*athena.Operator, members []*athena.Member, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("failed to marshal operators: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write(data)
	bs := h.Sum(nil)

	data, err = json.Marshal(members)
	if err != nil {
		return fmt.Errorf("Failed to marsahl payload: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(MEMBERS, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}
