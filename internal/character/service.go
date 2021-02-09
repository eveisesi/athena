package character

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Character(ctx context.Context, id uint, options ...OptionFunc) (*athena.Character, error)
	Characters(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Character, error)
	CharacterCorporationHistory(ctx context.Context, character *athena.Character, options ...OptionFunc) ([]*athena.CharacterCorporationHistory, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	corporation corporation.Service

	character athena.CharacterRepository
}

const (
	errPrefix = "[Character Service]"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, corporation corporation.Service, character athena.CharacterRepository) Service {
	return &service{
		logger: logger,

		cache:       cache,
		esi:         esi,
		corporation: corporation,

		character: character,
	}
}

func (s *service) Character(ctx context.Context, id uint, optionFuncs ...OptionFunc) (*athena.Character, error) {

	options := s.options(optionFuncs)

	etag, err := s.esi.Etag(ctx, esi.GetCharacter, esi.ModWithCharacter(&athena.Character{ID: id}))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	exists := true
	cached := true

	character, err := s.cache.Character(ctx, id)
	if err != nil {
		return nil, err
	}

	if character == nil {
		cached = false
		character, err = s.character.Character(ctx, character.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if character == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			character = &athena.Character{ID: character.ID}
			err = s.esi.ResetEtag(ctx, etag)
			if err != nil {
				return nil, err
			}
		}

	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetCharacter(ctx, character)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return character, nil

	}

	character, _, _, err = s.esi.GetCharacter(ctx, character)
	if err != nil {
		err = fmt.Errorf("[Character Service] Failed to fetch character %d from ESI: %w", id, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if options.history {
		s.CharacterCorporationHistory(ctx, character)
	}

	switch exists {
	case true:
		character, err = s.character.UpdateCharacter(ctx, character.ID, character)
		if err != nil {
			err = fmt.Errorf("[Character Service] Failed to create character %d in database: %w", id, err)
			return nil, err
		}
	case false:
		character, err = s.character.CreateCharacter(ctx, character)
		if err != nil {
			err = fmt.Errorf("[Character Service] Failed to create character %d in database: %w", id, err)
			return nil, err
		}

	}

	err = s.cache.SetCharacter(ctx, character)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return character, err

}

func (s *service) Characters(ctx context.Context, operators []*athena.Operator, options ...OptionFunc) ([]*athena.Character, error) {

	characters, err := s.character.Characters(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return characters, nil

}

func (s *service) CharacterCorporationHistory(ctx context.Context, character *athena.Character, optionFuncs ...OptionFunc) ([]*athena.CharacterCorporationHistory, error) {

	// options := s.options(optionFuncs)

	etag, err := s.esi.Etag(ctx, esi.GetCharacterCorporationHistory, esi.ModWithCharacter(character))
	if err != nil {
		return nil, err
	}

	exists := true
	cached := true

	history, err := s.cache.CharacterCorporationHistory(ctx, character.ID)
	if err != nil {
		return nil, err
	}

	if history == nil {
		cached = false
		character, err = s.character.Character(ctx, character.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			exists = false
			history = make([]*athena.CharacterCorporationHistory, 0)
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetCharacterCorporationHistory(ctx, character.ID, history)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return history, nil
	}

	newHistory, _, _, err := s.esi.GetCharacterCorporationHistory(ctx, character, make([]*athena.CharacterCorporationHistory, 0))
	if err != nil {
		return nil, fmt.Errorf("[Contacts Service] Failed to fetch corporation history for character %d: %w", character.ID, err)
	}

	if len(newHistory) > 0 {
		s.resolveHistoryAttributes(ctx, newHistory)
		history, err := s.diffAndUpdateHistory(ctx, character, history, newHistory)
		if err != nil {
			return nil, fmt.Errorf("%s Failed to diffAndUpdateContacts: %w", errPrefix, err)
		}

		err = s.cache.SetCharacterCorporationHistory(ctx, character.ID, history)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return history, err
}

func (s *service) resolveHistoryAttributes(ctx context.Context, history []*athena.CharacterCorporationHistory) {

	for _, record := range history {
		_, err := s.corporation.Corporation(ctx, record.CorporationID)
		if err != nil {
			s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
				"record_id":      record.RecordID,
				"corporation_id": record.CorporationID,
			}).Error("failed to resolve corporation record in character history")
		}
	}

}

func (s *service) diffAndUpdateHistory(ctx context.Context, character *athena.Character, old []*athena.CharacterCorporationHistory, new []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, error) {

	recordsToCreate := make([]*athena.CharacterCorporationHistory, 0)

	oldRecordMap := make(map[uint]*athena.CharacterCorporationHistory)
	for _, record := range old {
		oldRecordMap[record.RecordID] = record
	}

	for _, record := range new {
		if _, ok := oldRecordMap[record.RecordID]; !ok {
			recordsToCreate = append(recordsToCreate, record)
		}
	}

	var final = make([]*athena.CharacterCorporationHistory, 0)
	if len(recordsToCreate) > 0 {
		createdRecords, err := s.character.CreateCharacterCorporationHistory(ctx, character.ID, recordsToCreate)
		if err != nil {
			return nil, err
		}

		final = append(final, createdRecords...)
	}

	return final, nil

}
