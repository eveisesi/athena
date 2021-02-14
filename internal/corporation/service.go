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
	Corporation(ctx context.Context, id uint, options ...OptionFunc) (*athena.Corporation, error)
	Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	alliance alliance.Service

	corporation athena.CorporationRepository
}

const (
	errPrefix = "Corporation Service"
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

func (s *service) Corporation(ctx context.Context, id uint, optionFuncs ...OptionFunc) (*athena.Corporation, error) {

	options := s.options(optionFuncs)

	etag, err := s.esi.Etag(ctx, esi.GetCorporation, esi.ModWithCorporation(&athena.Corporation{ID: id}))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	exists := true
	cached := true

	corporation, err := s.cache.Corporation(ctx, id)
	if err != nil {
		return nil, err
	}

	if corporation == nil {
		cached = false
		corporation, err = s.corporation.Corporation(ctx, id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if corporation == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			corporation = &athena.Corporation{ID: id}
			err = s.esi.ResetEtag(ctx, etag)
			if err != nil {
				s.logger.WithError(err).WithField("id", id).
					Error("failed to reset for corporation alliance history")
			}
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetCorporation(ctx, corporation)
			if err != nil {
				s.logger.WithError(err).WithField("id", corporation.ID).
					Error("failed to fetch character corporation history")
			}
		}

		return corporation, nil
	}

	corporation, _, _, err = s.esi.GetCorporation(ctx, corporation)
	if err != nil {
		err = fmt.Errorf("[%s] Failed to fetch corporation %d from ESI: %w", errPrefix, corporation.ID, err)
		return nil, err
	}

	if options.history {
		_, err := s.CorporationAllianceHistory(ctx, corporation)
		if err != nil {
			s.logger.WithError(err).
				WithContext(ctx).WithField("id", corporation.ID).WithField("service", errPrefix).Error("Failed to fetch corporation history for corporation")
		}
	}

	switch exists {
	case true:
		corporation, err = s.corporation.UpdateCorporation(ctx, corporation.ID, corporation)
		if err != nil {
			err = fmt.Errorf("[%s] Failed to update corporation %d in the database: %w", errPrefix, corporation.ID, err)
			s.logger.WithError(err).WithContext(ctx).
				WithField("id", corporation.ID).Errorln()
			return nil, err
		}
	case false:
		corporation, err = s.corporation.CreateCorporation(ctx, corporation)
		if err != nil {
			err = fmt.Errorf("[%s] Failed to create corporation %d in the database: %w", errPrefix, corporation.ID, err)
			s.logger.WithError(err).WithContext(ctx).
				WithField("id", corporation.ID).Errorln()
			return nil, err
		}
	}

	err = s.cache.SetCorporation(ctx, corporation)
	if err != nil {
		err = fmt.Errorf("[%s] Failed to cache corporation %d: %w", errPrefix, corporation.ID, err)
		s.logger.WithError(err).WithContext(ctx).
			WithField("id", corporation.ID).Errorln()
		return nil, err
	}

	return corporation, err

}

func (s *service) Corporations(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Corporation, error) {

	corporations, err := s.corporation.Corporations(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return corporations, nil

}

func (s *service) CorporationAllianceHistory(ctx context.Context, corporation *athena.Corporation, optionFuncs ...OptionFunc) ([]*athena.CorporationAllianceHistory, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCorporationAllianceHistory, esi.ModWithCorporation(corporation))
	if err != nil {
		return nil, err
	}

	exists := true
	cached := true

	history, err := s.cache.CorporationAllianceHistory(ctx, corporation.ID)
	if err != nil {
		return nil, err
	}

	if history == nil {
		cached = false
		history, err := s.corporation.CorporationAllianceHistory(ctx, athena.NewEqualOperator("corporation_id", corporation.ID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		exists = history == nil || errors.Is(err, sql.ErrNoRows)

	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetCorporationAllianceHistory(ctx, corporation.ID, history)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}

		}

		return history, err

	}

	newHistory, _, _, err := s.esi.GetCorporationAllianceHistory(ctx, corporation, make([]*athena.CorporationAllianceHistory, 0))
	if err != nil {
		return nil, fmt.Errorf("[Contacts Service] Failed to fetch alliance history for corporation %d: %w", corporation.ID, err)
	}

	if len(newHistory) > 0 {
		s.resolveHistoryAttributes(ctx, newHistory)
		history, err := s.diffAndUpdateHistory(ctx, corporation, history, newHistory)
		if err != nil {
			return nil, fmt.Errorf("[%s] Failed to diffAndUpdateContacts: %w", errPrefix, err)
		}

		err = s.cache.SetCorporationAllianceHistory(ctx, corporation.ID, history)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return history, err

}

func (s *service) resolveHistoryAttributes(ctx context.Context, history []*athena.CorporationAllianceHistory) {

	for _, record := range history {
		if !record.AllianceID.Valid {
			continue
		}
		_, err := s.alliance.Alliance(ctx, record.AllianceID.Uint)
		if err != nil {
			s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
				"record_id":   record.RecordID,
				"alliance_id": record.AllianceID,
			}).Error("failed to resolve allilance record in corporation alliance history")
		}
	}

}
func (s *service) diffAndUpdateHistory(ctx context.Context, corporation *athena.Corporation, old []*athena.CorporationAllianceHistory, new []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, error) {

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
		createdRecords, err := s.corporation.CreateCorporationAllianceHistory(ctx, corporation.ID, recordsToCreate)
		if err != nil {
			return nil, err
		}

		final = append(final, createdRecords...)
	}

	return final, nil

}
