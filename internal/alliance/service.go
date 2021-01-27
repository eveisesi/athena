package alliance

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	AllianceByAllianceID(ctx context.Context, id uint, options []OptionFunc) (*athena.Alliance, error)
	Alliances(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Alliance, error)
	CreateAlliance(ctx context.Context, alliance *athena.Alliance, options []OptionFunc) (*athena.Alliance, error)
}

type service struct {
	cache cache.Service
	esi   esi.Service

	athena.AllianceRepository
}

func NewService(cache cache.Service, esi esi.Service, alliance athena.AllianceRepository) Service {
	return &service{
		cache: cache,
		esi:   esi,

		AllianceRepository: alliance,
	}
}

func (s *service) AllianceByAllianceID(ctx context.Context, id uint, options []OptionFunc) (*athena.Alliance, error) {

	alliances, err := s.Alliances(ctx, athena.NewOperators(athena.NewEqualOperator("alliance_id", id)), options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if len(alliances) == 1 {
		return alliances[0], nil
	}

	alliance, _, err := s.esi.GetAlliance(ctx, &athena.Alliance{AllianceID: id})
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	alliance, err = s.CreateAlliance(ctx, alliance, options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)

	}

	return alliance, nil

}

func (s *service) Alliances(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Alliance, error) {

	opts := s.options(options)

	if !opts.skipCache {
		alliances, err := s.cache.Alliances(ctx, operators)
		if err != nil {
			return nil, err
		}

		if alliances != nil {
			return alliances, nil
		}
	}

	alliances, err := s.AllianceRepository.Alliances(ctx, operators...)
	if err != nil {
		return nil, err
	}

	if opts.skipCache {
		return alliances, nil
	}

	err = s.cache.SetAlliances(ctx, operators, alliances, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return alliances, nil

}

func (s *service) CreateAlliance(ctx context.Context, alliance *athena.Alliance, options []OptionFunc) (*athena.Alliance, error) {

	alliance, err := s.AllianceRepository.CreateAlliance(ctx, alliance)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	opts := s.options(options)

	if opts.skipCache {
		return alliance, err
	}

	err = s.cache.SetAlliance(ctx, alliance, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return alliance, nil

}
