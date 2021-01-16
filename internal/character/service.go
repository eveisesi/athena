package character

import (
	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
)

type Service interface {
}

type service struct {
	cache       cache.Service
	esi         esi.Service
	alliance    alliance.Service
	corporation corporation.Service

	character athena.CharacterRepository
}

func NewService(cache cache.Service, esi esi.Service, alliance alliance.Service, corporation corporation.Service, character athena.CharacterRepository) Service {
	return &service{
		cache, esi, alliance, corporation, character,
	}
}
