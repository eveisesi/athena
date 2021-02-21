package corporation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FetchCorporation(ctx context.Context, corporationID uint) (*athena.Etag, error)
	Corporation(ctx context.Context, id uint) (*athena.Corporation, error)
	Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error)

	FetchCorporationAllianceHistory(ctx context.Context, corporationID uint) (*athena.Etag, error)
	CorporationAllianceHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CorporationAllianceHistory, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	alliance alliance.Service

	corporation athena.CorporationRepository
}

const (
	serviceIdentifier = "Corporation Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, alliance alliance.Service, corporation athena.CorporationRepository) Service {
	return &service{
		logger: logger,

		cache:    cache,
		esi:      esi,
		alliance: alliance,

		corporation: corporation,
	}
}

func (s *service) FetchCorporation(ctx context.Context, corporationID uint) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCorporation, esi.ModWithCorporationID(corporationID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"corporation_id": corporationID,
		"service":        serviceIdentifier,
		"method":         "FetchCorporation",
	})

	petag := etag.Etag
	corporation, etag, _, err := s.esi.GetCorporation(ctx, corporationID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch corporation from ESI")
		return nil, fmt.Errorf("failed to fetch corporation from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, err
	}

	existing, err := s.corporation.Corporation(ctx, corporationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch corporation from DB")
		return nil, fmt.Errorf("failed to fetch corporation from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		corporation, err = s.corporation.UpdateCorporation(ctx, corporationID, corporation)
		if err != nil {
			entry.WithError(err).Error("failed to update corporation in DB")
			return nil, fmt.Errorf("failed to update corporation in DB")
		}
	case false:
		corporation, err = s.corporation.CreateCorporation(ctx, corporation)
		if err != nil {
			entry.WithError(err).Error("failed to create corporation in DB")
			return nil, fmt.Errorf("failed to create corporation in DB")
		}
	}

	err = s.cache.SetCorporation(ctx, corporationID, corporation)
	if err != nil {
		entry.WithError(err).Error("failed to cache corporation")
		return nil, fmt.Errorf("failed to cache corporation")
	}

	return etag, err

}

func (s *service) Corporation(ctx context.Context, corporationID uint) (*athena.Corporation, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"alliance_id": corporationID,
		"service":     serviceIdentifier,
		"method":      "Corporation",
	})

	corporation, err := s.cache.Corporation(ctx, corporationID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch corporation from cache")
		return nil, fmt.Errorf("failed to fetch corporation from cache")
	}

	if corporation != nil {
		return corporation, nil
	}

	corporation, err = s.corporation.Corporation(ctx, corporationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch corporation from db")
		return nil, fmt.Errorf("failed to fetch corporation from db")
	}

	err = s.cache.SetCorporation(ctx, corporationID, corporation)
	if err != nil {
		entry.WithError(err).Error("failed to cache corporation")
	}

	return corporation, err

}

func (s *service) Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Alliances",
	})

	corporations, err := s.cache.Corporations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliances from cache")
		return nil, fmt.Errorf("failed to fetch alliances from cache")
	}

	if len(corporations) > 0 {
		return corporations, nil
	}

	corporations, err = s.corporation.Corporations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch alliances from db")
		return nil, fmt.Errorf("failed to fetch alliances from db")
	}

	err = s.cache.SetCorporations(ctx, corporations, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache corporations")
		return nil, fmt.Errorf("failed to cache corporations")
	}

	return corporations, nil

}

func (s *service) FetchCorporationAllianceHistory(ctx context.Context, corporationID uint) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCorporationAllianceHistory, esi.ModWithCorporationID(corporationID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"corporation_id": corporationID,
		"service":        serviceIdentifier,
		"method":         "FetchCorporationAllianceHistory",
	})

	petag := etag.Etag
	history, etag, _, err := s.esi.GetCorporationAllianceHistory(ctx, corporationID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch corporation alliance history from ESI")
		return nil, fmt.Errorf("failed to fetch corporation alliance history from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	s.resolveHistoryAttributes(ctx, history)

	existingHistory, err := s.CorporationAllianceHistory(ctx, athena.NewEqualOperator("corporation_id", corporationID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch existing history for corporation")
		return nil, fmt.Errorf("failed to fetch existing history for corporation")
	}

	_, err = s.diffAndUpdateHistory(ctx, corporationID, existingHistory, history)
	if err != nil {
		entry.WithError(err).Error("unexpected error encountered processing character history")
		return nil, fmt.Errorf("unexpected error encountered processing character history")
	}

	return etag, err

}

func (s *service) resolveHistoryAttributes(ctx context.Context, records []*athena.CorporationAllianceHistory) {
	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "resolveHistoryAttributes",
	})

	for _, record := range records {
		if !record.AllianceID.Valid {
			continue
		}
		_, err := s.alliance.FetchAlliance(ctx, record.AllianceID.Uint)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"record_id":      record.RecordID,
				"corporation_id": record.AllianceID.Uint,
			}).Error("failed to resolve alliance record in corporation history")
		}
	}

}

func (s *service) CorporationAllianceHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CorporationAllianceHistory, error) {

	history, err := s.cache.CorporationAllianceHistory(ctx, operators...)
	if err != nil {
		return nil, err
	}

	if history != nil {
		return history, nil
	}

	history, err = s.corporation.CorporationAllianceHistory(ctx, operators...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(history) > 0 {
		err = s.cache.SetCorporationAllianceHistory(ctx, history, operators...)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return history, err

}

func (s *service) diffAndUpdateHistory(ctx context.Context, corporationID uint, old []*athena.CorporationAllianceHistory, new []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, error) {

	recordsToCreate := make([]*athena.CorporationAllianceHistory, 0)

	oldRecordMap := make(map[uint]*athena.CorporationAllianceHistory)
	for _, record := range old {
		oldRecordMap[record.RecordID] = record
	}

	for _, record := range new {
		if _, ok := oldRecordMap[record.RecordID]; !ok {
			recordsToCreate = append(recordsToCreate, record)
		}
	}

	var final = make([]*athena.CorporationAllianceHistory, 0)
	if len(recordsToCreate) > 0 {
		createdRecords, err := s.corporation.CreateCorporationAllianceHistory(ctx, corporationID, recordsToCreate)
		if err != nil {
			return nil, err
		}

		final = append(final, createdRecords...)
	}

	return final, nil

}
