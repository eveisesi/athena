package athena

import (
	"context"

	"github.com/volatiletech/null"
)

type AuthRepository interface {
	AuthAttempt(ctx context.Context, hash string) (*AuthAttempt, error)
	CreateAuthAttempt(ctx context.Context, attempt *AuthAttempt) (*AuthAttempt, error)
	UpdateAuthAttempt(ctx context.Context, hash string, attempt *AuthAttempt) (*AuthAttempt, error)

	JSONWebKeySet(ctx context.Context) ([]byte, error)
	SaveJSONWebKeySet(ctx context.Context, jwks []byte) error
}

type AuthAttempt struct {
	Status AuthAttemptStatus
	State  string
	Token  null.String
	User   *Member
}

type AuthAttemptStatus string

const (
	CreateAuthStatus    AuthAttemptStatus = "created"
	PendingAuthStatus   AuthAttemptStatus = "pending"
	InvalidAuthStatus   AuthAttemptStatus = "invalid"
	ExpiredAuthStatus   AuthAttemptStatus = "expired"
	CompletedAuthStatus AuthAttemptStatus = "completed"
)

func (a AuthAttemptStatus) String() string {
	return string(a)
}
