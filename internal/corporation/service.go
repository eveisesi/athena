package corporation

import (
	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
)

type Service interface {
}

type service struct {
	cache cache.Service
	esi   esi.Service

	corporation athena.CorporationRepository
}

func NewService(cache cache.Service, esi esi.Service, corporation athena.CorporationRepository) Service {
	return &service{
		cache, esi, corporation,
	}
}
