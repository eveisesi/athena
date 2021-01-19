package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type Service interface {
	InitializeAttempt(ctx context.Context) (*athena.AuthAttempt, error)
	AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error)
	UpdateAuthAttempt(ctx context.Context, hash string, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error)

	ValidateToken(ctx context.Context, member *athena.Member) (*athena.Member, error)
	AuthorizationURI(ctx context.Context, state string) string
	BearerForCode(ctx context.Context, code string) (*oauth2.Token, error)
	ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error)
}

type service struct {
	// athena.AuthRepository
	oauth *oauth2.Config
	cache cache.Service

	client  *http.Client
	jwksURI string
}

// authRepo athena.AuthRepository
func NewService(cache cache.Service, oauth *oauth2.Config, client *http.Client, jwksURI string) *service {
	return &service{
		cache:   cache,
		oauth:   oauth,
		client:  client,
		jwksURI: jwksURI,
	}
}

func (s *service) InitializeAttempt(ctx context.Context) (*athena.AuthAttempt, error) {

	h := hmac.New(sha256.New, nil)
	_, _ = h.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	b := h.Sum(nil)

	attempt := &athena.AuthAttempt{
		Status: athena.PendingAuthStatus,
		State:  fmt.Sprintf("%x", string(b)),
	}

	return s.cache.CreateAuthAttempt(ctx, attempt)

}

func (s *service) AuthAttempt(ctx context.Context, hash string) (*athena.AuthAttempt, error) {

	attempt, err := s.cache.AuthAttempt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attempt with hash of %s: %w", hash, err)
	}

	if attempt == nil {
		attempt = new(athena.AuthAttempt)
		attempt.Status = athena.InvalidAuthStatus
	}

	return attempt, nil

}

func (s *service) UpdateAuthAttempt(ctx context.Context, hash string, attempt *athena.AuthAttempt) (*athena.AuthAttempt, error) {

	attempt.State = hash

	attempt, err := s.cache.CreateAuthAttempt(ctx, attempt)
	if err != nil {
		return nil, err
	}

	return attempt, nil

}

func (s *service) AuthorizationURI(ctx context.Context, state string) string {
	return s.oauth.AuthCodeURL(state)
}

func (s *service) ValidateToken(ctx context.Context, member *athena.Member) (*athena.Member, error) {

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)

	token := new(oauth2.Token)
	token.AccessToken = member.AccessToken
	token.RefreshToken = member.RefreshToken
	token.Expiry = member.Expires

	tokenSource := s.oauth.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if member.AccessToken != newToken.AccessToken {
		member.AccessToken = newToken.AccessToken
		member.Expires = newToken.Expiry
		member.RefreshToken = newToken.RefreshToken
	}

	return member, nil
}

func (s *service) BearerForCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth.Exchange(ctx, code)
}

func (s *service) ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error) {

	set, err := s.getSet()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jwks: %w", err)
	}

	token, err := jwt.ParseString(t, jwt.WithKeySet(set))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil

}

func (s *service) getSet() (*jwk.Set, error) {

	ctx := context.Background()

	b, err := s.cache.JSONWebKeySet(ctx)
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured querying redis for jwks: %w", err)
	}

	if b == nil {
		res, err := s.client.Get(s.jwksURI)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve jwks from sso: %w", err)
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code recieved while fetching jwks. %d", res.StatusCode)
		}

		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read jwk response body: %w", err)
		}

		err = s.cache.SaveJSONWebKeySet(ctx, buf, cache.ExpiryHours(6))
		if err != nil {
			return nil, fmt.Errorf("failed to save jwks to cache layer: %w", err)
		}

		b = buf
	}

	return jwk.ParseBytes(b)

}
