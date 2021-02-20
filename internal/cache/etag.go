package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type etagService interface {
	Etag(ctx context.Context, etagID string) (*athena.Etag, error)
	SetEtag(ctx context.Context, etagID string, etag *athena.Etag, expires time.Duration) error
	DeleteEtag(ctx context.Context, etagID string) error
}

const (
	keyEtag = "athena::etag::%s"
)

func (s *service) Etag(ctx context.Context, etagID string) (*athena.Etag, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyEtag, etagID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var etag = new(athena.Etag)

	err = json.Unmarshal([]byte(result), &etag)
	if err != nil {
		return nil, err
	}

	return etag, nil

}
func (s *service) SetEtag(ctx context.Context, etagID string, etag *athena.Etag, expires time.Duration) error {

	data, err := json.Marshal(etag)
	if err != nil {
		return fmt.Errorf("failed to marshal etag: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyEtag, etagID), data, expires).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) DeleteEtag(ctx context.Context, etagID string) error {

	key := fmt.Sprintf(keyEtag, etagID)
	_, err := s.client.Del(ctx, fmt.Sprintf(keyEtag, etagID)).Result()
	if err != nil {
		return fmt.Errorf("[Cache Service] Failed to delete key %s from cache: %w", key, err)
	}

	return nil

}
