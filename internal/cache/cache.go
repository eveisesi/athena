package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Service interface {
	authService
	memberService
}

type service struct {
	client *redis.Client
}

func NewService(client *redis.Client) Service {
	return &service{
		client: client,
	}
}

type options struct {
	expiry time.Duration
}

type OptionsFunc func(opts *options) *options

func defaultOptions() *options {
	return &options{
		expiry: time.Minute * 5,
	}
}

func applyOptionFuncs(opts *options, optionFuncs []OptionsFunc) *options {
	if opts == nil {
		opts = defaultOptions()
	}

	for _, optionFunc := range optionFuncs {
		opts = optionFunc(opts)
	}

	return opts
}

func WithCustomExpiry(expiry time.Duration) OptionsFunc {
	return func(opts *options) *options {
		opts.expiry = expiry
		return opts
	}
}
