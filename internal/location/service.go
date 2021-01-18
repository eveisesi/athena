package location

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/sirupsen/logrus"
)

type Service interface {
	MemberLocation(ctx context.Context, member *athena.Member) error
	MemberShip(ctx context.Context, member *athena.Member) error
	MemberOnline(ctx context.Context, member *athena.Member) error
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	location athena.MemberLocationRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, location athena.MemberLocationRepository) Service {
	return &service{
		logger: logger,

		cache: cache,
		esi:   esi,

		location: location,
	}
}

func (s *service) MemberLocation(ctx context.Context, member *athena.Member) error {
	return nil
}

func (s *service) MemberShip(ctx context.Context, member *athena.Member) error {
	return nil
}

func (s *service) MemberOnline(ctx context.Context, member *athena.Member) error {
	return nil
}
