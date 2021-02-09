package cache

import (
	"time"

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
	memberService
	processorService
	skillService
	universeService
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

type options struct {
	expiry time.Duration
}

type OptionFunc func(opts *options) *options

func defaultOptions() *options {
	return &options{
		expiry: time.Minute * 5,
	}
}

func applyOptionFuncs(opts *options, optionFuncs []OptionFunc) *options {
	if opts == nil {
		opts = defaultOptions()
	}

	for _, optionFunc := range optionFuncs {
		opts = optionFunc(opts)
	}

	return opts
}

func WithCustomExpiry(expiry time.Duration) OptionFunc {
	return func(opts *options) *options {
		opts.expiry = expiry
		return opts
	}
}

func ExpiryMinutes(min int) OptionFunc {
	return WithCustomExpiry(time.Minute * time.Duration(min))
}

func ExpiryHours(hr int) OptionFunc {
	return WithCustomExpiry(time.Hour * time.Duration(hr))
}
