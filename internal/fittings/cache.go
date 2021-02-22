package fittings

import "github.com/go-redis/redis/v8"

type cache interface {
	MemberFittings()
}

type cacher struct {
	client *redis.Client
}

func newCacher(client *redis.Client) cache {
	return &cacher{
		client: client,
	}
}

func (s *cacher) MemberFittings() {}
