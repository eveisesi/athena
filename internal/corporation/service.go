package corporation

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	CorporationByCorporationID(ctx context.Context, id uint, options []OptionFunc) (*athena.Corporation, error)
	Corporations(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Corporation, error)
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

func (s *service) CorporationByCorporationID(ctx context.Context, id uint, options []OptionFunc) (*athena.Corporation, error) {

	corporations, err := s.Corporations(ctx, athena.NewOperators(athena.NewEqualOperator("corporation_id", id)), options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if len(corporations) == 1 {
		return corporations[0], nil
	}

	corporation, _, err := s.esi.GetCorporation(ctx, &athena.Corporation{CorporationID: id})
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

func (s *service) Corporations(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Corporation, error) {

	opts := s.options(options)

	if !opts.skipCache {
		corporations, err := s.cache.Corporations(ctx, operators)
		if err != nil {
			return nil, err
		}

		if corporations != nil {
			return corporations, nil
		}
	}

	corporations, err := s.CorporationRepository.Corporations(ctx, operators...)
	if err != nil {
		return nil, err
	}

	if opts.skipCache {
		return corporations, nil
	}

	err = s.cache.SetCorporations(ctx, operators, corporations, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return corporations, nil

}

func (s *service) CreateCorporation(ctx context.Context, corporation *athena.Corporation, options []OptionFunc) (*athena.Corporation, error) {

	corporation, err := s.CorporationRepository.CreateCorporation(ctx, corporation)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	opts := s.options(options)

	if opts.skipCache {
		return corporation, err
	}

	err = s.cache.SetCorporation(ctx, corporation, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return corporation, nil

}
