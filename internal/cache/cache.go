package cache

import (
	"github.com/go-redis/redis/v8"
)

type Service interface {
	allianceService
	authService
	characterService
	cloneService
	contactService
	contractService
	corporationService
	esiService
	etagService
	locationService
	mailService
	memberService
	processorService
	skillService
	universeService
	walletService
}

type service struct {
	client *redis.Client
}

const (
	errKeyNotFound        = "[Cache Service] Failed to fetch set member for key %s: %w"
	errSetUnmarshalFailed = "[Cache Service] Failed to unmarshal set member for key %s onto struct: %w"
)

// NewService returns a new instance of the Cache Service
func NewService(client *redis.Client) Service {
	return &service{
		client: client,
	}
}
