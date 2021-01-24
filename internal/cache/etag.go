package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type etagService interface {
	Etag(ctx context.Context, etagID string) (*athena.Etag, error)
	SetEtag(ctx context.Context, etagID string, etag *athena.Etag, optionFunc ...OptionFunc) error
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
func (s *service) SetEtag(ctx context.Context, etagID string, etag *athena.Etag, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(etag)
	if err != nil {
		return fmt.Errorf("failed to marshal etag: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCharacter, etagID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}