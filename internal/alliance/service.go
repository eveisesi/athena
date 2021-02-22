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
	FetchAlliance(ctx context.Context, allianceID uint) (*athena.Etag, error)
	Alliance(ctx context.Context, allianceID uint) (*athena.Alliance, error)
	Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error)
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

func (s *service) FetchAlliance(ctx context.Context, allianceID uint) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetAlliance, esi.ModWithAllianceID(allianceID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": allianceID,
		"service":   serviceIdentifier,
		"method":    "FetchAlliance",
	})

	alliance, etag, _, err := s.esi.GetAlliance(ctx, allianceID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliance from ESI")
		return nil, fmt.Errorf("failed to fetch alliance from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, err
	}

	existing, err := s.alliance.Alliance(ctx, allianceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch alliance from DB")
		return nil, fmt.Errorf("failed to fetch alliance from DB")
	}

	switch existing != nil {
	case true:
		alliance, err = s.alliance.UpdateAlliance(ctx, allianceID, alliance)
		if err != nil {
			entry.WithError(err).Error("failed to update alliance in DB")
			return nil, fmt.Errorf("failed to update alliance in DB")
		}
	case false:
		alliance, err = s.alliance.CreateAlliance(ctx, alliance)
		if err != nil {
			entry.WithError(err).Error("failed to create alliance in DB")
			return nil, fmt.Errorf("failed to create alliance in DB")
		}
	}

	err = s.cache.SetAlliance(ctx, allianceID, alliance)
	if err != nil {
		entry.WithError(err).Error("failed to cache alliance in Redis")
		return nil, fmt.Errorf("failed to cache alliance in Redis")
	}

	return etag, nil

}

func (s *service) Alliance(ctx context.Context, allianceID uint) (*athena.Alliance, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"alliance_id": allianceID,
		"service":     serviceIdentifier,
		"method":      "Alliance",
	})

	alliance, err := s.cache.Alliance(ctx, allianceID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliance from cache")
		return nil, fmt.Errorf("failed to fetch alliance from cache")
	}

	if alliance != nil {
		return alliance, nil
	}

	alliance, err = s.alliance.Alliance(ctx, allianceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch alliance from DB")
		return nil, fmt.Errorf("failed to fetch alliance from DB")
	}

	err = s.cache.SetAlliance(ctx, alliance.ID, alliance)
	if err != nil {
		entry.WithError(err).Error("failed to cache alliance")
	}

	return alliance, err

}

func (s *service) Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Alliances",
	})

	alliances, err := s.cache.Alliances(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliances from cache")
		return nil, fmt.Errorf("failed to fetch alliances from cache")
	}

	if len(alliances) > 0 {
		return alliances, nil
	}

	alliances, err = s.alliance.Alliances(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliances from db")
		return nil, fmt.Errorf("failed to fetch alliances from db")
	}

	err = s.cache.SetAlliances(ctx, alliances, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache alliances")
		return nil, fmt.Errorf("failed to cache alliances")
	}

	return alliances, nil

}
