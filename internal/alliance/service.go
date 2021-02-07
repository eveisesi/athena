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
	Alliance(ctx context.Context, id uint, options []OptionFunc) (*athena.Alliance, error)
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

func (s *service) Alliance(ctx context.Context, id uint, options []OptionFunc) (*athena.Alliance, error) {

	alliance, err := s.AllianceRepository.Alliance(ctx, id)
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

	alliance, err = s.CreateAlliance(ctx, alliance, options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return alliance, nil

}

func (s *service) Alliances(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Alliance, error) {

	alliances, err := s.AllianceRepository.Alliances(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return alliances, nil

}

func (s *service) CreateAlliance(ctx context.Context, alliance *athena.Alliance, options []OptionFunc) (*athena.Alliance, error) {

	// err := s.AllianceRepository.CreateAlliance(ctx, alliance)
	// if err != nil {
	// 	newrelic.FromContext(ctx).NoticeError(err)
	// 	return nil, err
	// }

	return nil, nil

}
