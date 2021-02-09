package alliance

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
	Alliance(ctx context.Context, id uint, options ...OptionFunc) (*athena.Alliance, error)
	Alliances(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Alliance, error)
}

type service struct {
	cache cache.Service
	esi   esi.Service

	alliance athena.AllianceRepository
}

func NewService(cache cache.Service, esi esi.Service, alliance athena.AllianceRepository) Service {
	return &service{
		cache: cache,
		esi:   esi,

		alliance: alliance,
	}
}

func (s *service) Alliance(ctx context.Context, id uint, options ...OptionFunc) (*athena.Alliance, error) {

	alliance, err := s.alliance.Alliance(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if err == nil && alliance != nil {
		return alliance, err
	}

	alliance, _, err = s.esi.GetAlliance(ctx, &athena.Alliance{ID: id})
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	alliance, err = s.alliance.CreateAlliance(ctx, alliance)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	err = s.cache.SetAlliance(ctx, alliance)

	return alliance, err

}

func (s *service) Alliances(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Alliance, error) {

	alliances, err := s.alliance.Alliances(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return alliances, nil

}
