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
	"github.com/sirupsen/logrus"
)

type Service interface {
	FetchCharacter(ctx context.Context, characterID uint) (*athena.Etag, error)
	Character(ctx context.Context, characterID uint) (*athena.Character, error)
	Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error)

	FetchCharacterCorporationHistory(ctx context.Context, characterID uint) (*athena.Etag, error)
	CharacterCorporationHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CharacterCorporationHistory, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	corporation corporation.Service

	character athena.CharacterRepository
}

const (
	serviceIdentifier = "Character Service"
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

func (s *service) FetchCharacter(ctx context.Context, characterID uint) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacter, esi.ModWithCharacterID(characterID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": characterID,
		"service":   serviceIdentifier,
		"method":    "FetchCharacter",
	})

	ptag := etag.Etag
	character, etag, _, err := s.esi.GetCharacter(ctx, characterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character from ESI")
		return nil, fmt.Errorf("failed to fetch character from ESI")
	}

	if etag.Etag == ptag {
		return etag, nil
	}

	existing, err := s.character.Character(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch character from DB")
		return nil, fmt.Errorf("failed to fetch character from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		character, err = s.character.UpdateCharacter(ctx, characterID, character)
		if err != nil {
			entry.WithError(err).Error("failed to update character in DB")
			return nil, fmt.Errorf("failed to update character in DB")
		}
	case false:
		character, err = s.character.CreateCharacter(ctx, character)
		if err != nil {
			entry.WithError(err).Error("failed to create character in DB")
			return nil, fmt.Errorf("failed to create character in DB")
		}
	}

	err = s.cache.SetCharacter(ctx, characterID, character)
	if err != nil {
		entry.WithError(err).Error("failed to cache character")
	}

	return etag, nil

}

func (s *service) FetchCharacterCorporationHistory(ctx context.Context, characterID uint) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterCorporationHistory, esi.ModWithCharacterID(characterID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"character_id": characterID,
		"service":      serviceIdentifier,
		"method":       "FetchCharacterCorporationHistory",
	})

	petag := etag.Etag
	history, etag, _, err := s.esi.GetCharacterCorporationHistory(ctx, characterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character from ESI")
		return nil, fmt.Errorf("failed to fetch character from ESI")
	}

	if etag.Etag == petag {
		return etag, nil
	}

	s.resolveHistoryAttributes(ctx, history)

	existingHistory, err := s.CharacterCorporationHistory(ctx, athena.NewEqualOperator("character_id", characterID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch existing history for character")
		return nil, fmt.Errorf("failed to fetch existing history for character")
	}

	_, err = s.diffAndUpdateHistory(ctx, characterID, existingHistory, history)
	if err != nil {
		entry.WithError(err).Error("unexpected error encountered processing character history")
		return nil, fmt.Errorf("unexpected error encountered processing character history")
	}

	return etag, err

}

func (s *service) Character(ctx context.Context, characterID uint) (*athena.Character, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"character_id": characterID,
		"service":      serviceIdentifier,
		"method":       "Character",
	})

	character, err := s.cache.Character(ctx, characterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character from cache")
		return nil, fmt.Errorf("failed to fetch character from cache")
	}

	if character != nil {
		return character, nil
	}

	character, err = s.character.Character(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch character from DB")
		return nil, fmt.Errorf("failed to fetch character from DB")
	}

	err = s.cache.SetCharacter(ctx, characterID, character)
	if err != nil {
		entry.WithError(err).Error("failed to cache character")
	}

	return character, err

}

func (s *service) Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Characters",
	})

	characters, err := s.cache.Characters(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch characters from cache")
		return nil, fmt.Errorf("failed to fetch characters from cache")
	}

	if len(characters) > 0 {
		return characters, nil
	}

	characters, err = s.character.Characters(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch characters from db")
		return nil, fmt.Errorf("failed to fetch characters from db")
	}

	err = s.cache.SetCharacters(ctx, characters, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache characters")
	}

	return characters, nil

}

func (s *service) CharacterCorporationHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CharacterCorporationHistory, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Characters",
	})

	history, err := s.cache.CharacterCorporationHistory(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character corporation history from cache")
		return nil, fmt.Errorf("failed to fetch character corporation history from cache")
	}

	if len(history) > 0 {
		return history, nil
	}

	history, err = s.character.CharacterCorporationHistory(ctx, operators...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch character corporation from db")
		return nil, fmt.Errorf("failed to fetch character corporation from db")
	}

	if len(history) == 0 {
		return nil, nil
	}

	err = s.cache.SetCharacterCorporationHistory(ctx, history, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache character corporation history")
	}

	return history, nil

}

func (s *service) resolveHistoryAttributes(ctx context.Context, history []*athena.CharacterCorporationHistory) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "resolveHistoryAttributes",
	})

	for _, record := range history {
		_, err := s.corporation.Corporation(ctx, record.CorporationID)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"record_id":      record.RecordID,
				"corporation_id": record.CorporationID,
			}).Error("failed to resolve corporation record in character history")
		}
	}

}

func (s *service) diffAndUpdateHistory(ctx context.Context, characterID uint, old []*athena.CharacterCorporationHistory, new []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"character_id": characterID,
		"service":      serviceIdentifier,
		"method":       "diffAndUpdateHistory",
	})

	recordsToCreate := make([]*athena.CharacterCorporationHistory, 0)

	oldRecordMap := make(map[uint64]*athena.CharacterCorporationHistory)
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
		createdRecords, err := s.character.CreateCharacterCorporationHistory(ctx, characterID, recordsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create character corporation history records in db")
		}

		final = append(final, createdRecords...)
	}

	return final, nil

}
