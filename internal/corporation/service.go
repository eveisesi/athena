package corporation

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	Corporation(ctx context.Context, id uint, options []OptionFunc) (*athena.Corporation, error)
	Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error)
	CreateCorporation(ctx context.Context, corporation *athena.Corporation, options []OptionFunc) (*athena.Corporation, error)
}

type service struct {
	cache cache.Service
	esi   esi.Service

	athena.CorporationRepository
}

func NewService(cache cache.Service, esi esi.Service, corporation athena.CorporationRepository) Service {
	return &service{
		cache: cache,
		esi:   esi,

		CorporationRepository: corporation,
	}
}

func (s *service) Corporation(ctx context.Context, id uint, options []OptionFunc) (*athena.Corporation, error) {

	corporation, err := s.CorporationRepository.Corporation(ctx, id)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	corporation, _, err = s.esi.GetCorporation(ctx, &athena.Corporation{CorporationID: id})
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	corporation, err = s.CreateCorporation(ctx, corporation, options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return corporation, nil

}

func (s *service) Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error) {

	corporations, err := s.CorporationRepository.Corporations(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return corporations, nil

}

func (s *service) CreateCorporation(ctx context.Context, corporation *athena.Corporation, options []OptionFunc) (*athena.Corporation, error) {
	return s.CorporationRepository.CreateCorporation(ctx, corporation)
}
