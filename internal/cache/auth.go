package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ulule/deepcopier"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type authRepository struct {
	client *redis.Client
	expiry time.Duration
}

const AUTH_ATTEMPT = "athena::auth::attempt::%s"
const AUTH_WEB_KEY_SET = "athena::auth::jwks"

// NewAuthRepository returns an instance of authRepository that satisfies the athena.AuthRepository
// interface.
func NewAuthRepository(client *redis.Client, expiry time.Duration) athena.AuthRepository {

	return &authRepository{
		client,
		expiry,
	}

}

func (r *authRepository) JSONWebKeySet(ctx context.Context) ([]byte, error) {

	result, err := r.client.Get(ctx, AUTH_WEB_KEY_SET).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil

}

func (r *authRepository) SaveJSONWebKeySet(ctx context.Context, jwks []byte) error {

	_, err := r.client.Set(ctx, AUTH_WEB_KEY_SET, jwks, time.Hour*24).Result()

	return err

}

func (r *authRepository) AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error) {

	var attempt = new(athena.AuthAttempt)

	result, err := r.client.Get(ctx, fmt.Sprintf(AUTH_ATTEMPT, hash)).Bytes()
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

func (r *authRepository) CreateAuthAttempt(ctx context.Context, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error) {

	if attempt.State == "" {
		return nil, fmt.Errorf("empty state provided")
	}

	b, err := json.Marshal(attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	_, err = r.client.Set(ctx, fmt.Sprintf(AUTH_ATTEMPT, attempt.State), b, time.Minute*5).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth attempt: %w", err)
	}

	return attempt, nil
}

func (r *authRepository) UpdateAuthAttempt(ctx context.Context, hash string, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error) {

	original, err := r.AuthAttempt(ctx, hash)
	if err != nil {
		return nil, err
	}

	if original.Status == athena.InvalidAuthStatus {
		return original, nil
	}

	err = deepcopier.Copy(attempt).To(original)
	if err != nil {
		return nil, fmt.Errorf("failed to update auth attempt: %w", err)
	}

	original.State = hash

	b, err := json.Marshal(original)
	if err != nil {
		return nil, fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	_, err = r.client.Set(ctx, fmt.Sprintf(AUTH_ATTEMPT, original.State), b, time.Minute*5).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth attempt: %w", err)
	}

	return original, nil

}
