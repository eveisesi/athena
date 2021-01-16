package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type authService interface {
	JSONWebKeySet(ctx context.Context) ([]byte, error)
	SaveJSONWebKeySet(ctx context.Context, jwks []byte, optionFuncs ...OptionsFunc) error
	AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error)
	CreateAuthAttempt(ctx context.Context, attempt *athena.AuthAttempt, optionFuncs ...OptionsFunc) (*athena.AuthAttempt, error)
}

const AUTH_ATTEMPT = "athena::auth::attempt::%s"
const AUTH_WEB_KEY_SET = "athena::auth::jwks"

func (s *service) JSONWebKeySet(ctx context.Context) ([]byte, error) {

	result, err := s.client.Get(ctx, AUTH_WEB_KEY_SET).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil

}

func (s *service) SaveJSONWebKeySet(ctx context.Context, jwks []byte, optionFuncs ...OptionsFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	_, err := s.client.Set(ctx, AUTH_WEB_KEY_SET, jwks, options.expiry).Result()

	return err

}

func (s *service) AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error) {

	var attempt = new(athena.AuthAttempt)

	result, err := s.client.Get(ctx, fmt.Sprintf(AUTH_ATTEMPT, hash)).Bytes()
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

func (s *service) CreateAuthAttempt(ctx context.Context, attempt *athena.AuthAttempt, optionFuncs ...OptionsFunc) (*athena.AuthAttempt, error) {

	if attempt.State == "" {
		return nil, fmt.Errorf("empty state provided")
	}

	options := applyOptionFuncs(nil, optionFuncs)

	b, err := json.Marshal(attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(AUTH_ATTEMPT, attempt.State), b, options.expiry).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth attempt: %w", err)
	}

	return attempt, nil
}
