package alliance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Alliance(ctx context.Context, id uint, options ...OptionFunc) (*athena.Alliance, error)
	Alliances(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Alliance, error)
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	alliance athena.AllianceRepository
}

const (
	serviceIdentifier = "Alliance Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, alliance athena.AllianceRepository) Service {
	return &service{
		logger: logger,

		cache: cache,
		esi:   esi,

		alliance: alliance,
	}
}

func (s *service) Alliance(ctx context.Context, id uint, optionFuncs ...OptionFunc) (*athena.Alliance, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"alliance_id": id,
		"service":     serviceIdentifier,
		"method":      "Alliance",
	})

	etag, err := s.esi.Etag(ctx, esi.GetAlliance, esi.ModWithAlliance(&athena.Alliance{ID: id}))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	exists := true
	cached := true

	alliance, err := s.cache.Alliance(ctx, id)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliance from cache")
		return nil, fmt.Errorf("failed to fetch alliance from cache")
	}

	if alliance == nil {
		cached = false
		alliance, err = s.alliance.Alliance(ctx, id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to fetch alliance from DB")
			return nil, fmt.Errorf("failed to fetch alliance from DB")
		}

		if alliance == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			alliance = &athena.Alliance{ID: id}
			err = s.esi.ResetEtag(ctx, etag)
			if err != nil {
				entry.WithError(err).Error("failed to reset etag")
				return nil, fmt.Errorf("failed to reset etag")
			}
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetAlliance(ctx, alliance)
			if err != nil {
				entry.WithError(err).Error("failed to cache alliance")
			}
		}

		return alliance, nil

	}

	alliance, _, _, err = s.esi.GetAlliance(ctx, alliance)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliance from ESI")
		return nil, fmt.Errorf("failed to fetch alliance from ESI")
	}

	switch exists {
	case true:
		alliance, err = s.alliance.UpdateAlliance(ctx, alliance.ID, alliance)
		if err != nil {
			entry.WithError(err).Error("failed to update alliance in database")
			return nil, err
		}
	case false:
		alliance, err = s.alliance.CreateAlliance(ctx, alliance)
		if err != nil {
			entry.WithError(err).Error("failed to create alliance in database")
			return nil, err
		}

	}

	err = s.cache.SetAlliance(ctx, alliance)
	if err != nil {
		entry.WithError(err).Error("failed to cache alliance")
	}

	return alliance, err

}

func (s *service) Alliances(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Alliance, error) {

	alliances, err := s.alliance.Alliances(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return alliances, nil

}
