package corporation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	Corporation(ctx context.Context, id uint, options ...OptionFunc) (*athena.Corporation, error)
	Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error)
}

type service struct {
	cache cache.Service
	esi   esi.Service

	corporation athena.CorporationRepository
}

func NewService(cache cache.Service, esi esi.Service, corporation athena.CorporationRepository) Service {
	return &service{
		cache: cache,
		esi:   esi,

		corporation: corporation,
	}
}

func (s *service) Corporation(ctx context.Context, id uint, optionFunc ...OptionFunc) (*athena.Corporation, error) {

	corporation, err := s.corporation.Corporation(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if err == nil && corporation != nil {
		return corporation, err
	}

	corporation, _, err = s.esi.GetCorporation(ctx, &athena.Corporation{ID: id})
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	corporation, err = s.corporation.CreateCorporation(ctx, corporation)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	err = s.cache.SetCorporation(ctx, corporation)

	return corporation, err

}

func (s *service) Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error) {

	corporations, err := s.corporation.Corporations(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return corporations, nil

}
