package alliance

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

	alliance athena.AllianceRepository
}

func NewService(cache cache.Service, esi esi.Service, alliance athena.AllianceRepository) Service {
	return &service{
		cache, esi, alliance,
	}
}
