package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type authService interface {
	JSONWebKeySet(ctx context.Context) ([]byte, error)
	SaveJSONWebKeySet(ctx context.Context, jwks []byte) error
	AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error)
	CreateAuthAttempt(ctx context.Context, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error)
}

const keyAuthAttempt = "athena::auth::attempt::%s"
const keyAuthJWKS = "athena::auth::jwks"

func (s *service) JSONWebKeySet(ctx context.Context) ([]byte, error) {

	result, err := s.client.Get(ctx, keyAuthJWKS).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil

}

func (s *service) SaveJSONWebKeySet(ctx context.Context, jwks []byte) error {

	_, err := s.client.Set(ctx, keyAuthJWKS, jwks, time.Hour*6).Result()

	return err

}

func (s *service) AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error) {

	var attempt = new(athena.AuthAttempt)

	result, err := s.client.Get(ctx, fmt.Sprintf(keyAuthAttempt, hash)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(result, attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data onto result struct: %w", err)
	}

	return attempt, nil

}

func (s *service) CreateAuthAttempt(ctx context.Context, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error) {

	if attempt.State == "" {
		return nil, fmt.Errorf("empty state provided")
	}

	b, err := json.Marshal(attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyAuthAttempt, attempt.State), b, time.Minute*5).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth attempt: %w", err)
	}

	return attempt, nil
}
