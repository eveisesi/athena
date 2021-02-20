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
	Character(ctx context.Context, id uint) (*athena.Character, error)
	Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error)
	// CharacterCorporationHistory(ctx context.Context, operators []*athena.Operator) ([]*athena.CharacterCorporationHistory, error)
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

	character, etag, _, err := s.esi.GetCharacter(ctx, characterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character from ESI")
		return nil, fmt.Errorf("failed to fetch character from ESI")
	}

	if character == nil {
		return etag, err
	}

	existing, err := s.character.Character(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch character from DB")
		return nil, fmt.Errorf("failed to fetch character from DB")
	}

	exists := true

	if existing == nil || errors.Is(err, sql.ErrNoRows) {
		exists = false
	}

	switch exists {
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

	err = s.cache.SetCharacter(ctx, character, cache.ExpiryHours(1))
	if err != nil {
		entry.WithError(err).Error("failed to cache character in Redis")
		return nil, fmt.Errorf("failed to cache character in Redis")
	}

	return etag, nil

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

	if character == nil {
		character, err = s.character.Character(ctx, characterID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to fetch character from DB")
			return nil, fmt.Errorf("failed to fetch character from DB")
		}

		if character != nil {
			err = s.cache.SetCharacter(ctx, character)
			if err != nil {
				entry.WithError(err).Error("failed to cache character")
			}
		}
	}

	return character, err

}

func (s *service) Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Characters",
	})

	// characters, err := s.cache.Characters(ctx, operators)
	// if err != nil {
	// 	entry.WithError(err).Error("failed to fetch characters from cache")
	// 	return nil, fmt.Errorf("failed to fetch characters from cache")
	// }

	// if len(characters) > 0 {
	// 	return characters, nil
	// }

	characters, err := s.character.Characters(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch characters from cache")
		return nil, fmt.Errorf("failed to fetch characters from cache")
	}

	// err = s.cache.SetCharacters(ctx, operators, characters, cache.ExpiryMinutes(5))
	// if err != nil {
	// 	entry.WithError(err).Error("failed to cache characters in Redis")
	// 	return nil, fmt.Errorf("failed to cache characters in Redis")
	// }

	return characters, nil

}

// func (s *service) FetchCharacter(ctx context.Context, characterID uint) (*athena.Etag, error) {

// 	etag, err := s.esi.Etag(ctx, esi.GetCharacter, esi.ModWithCharacter(&athena.Character{ID: characterID}))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
// 	}

// }

// func (s *service) Character(ctx context.Context, id uint) (*athena.Character, error) {

// 	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
// 		"character_id": id,
// 		"service":      serviceIdentifier,
// 		"method":       "Character",
// 	})

// 	etag, err := s.esi.Etag(ctx, esi.GetCharacter, esi.ModWithCharacter(&athena.Character{ID: id}))
// 	if err != nil {
// 		entry.WithError(err).Error("failed to fetch etag object")
// 		return nil, fmt.Errorf("failed to fetch etag object")
// 	}

// 	exists := true
// 	cached := true

// 	character, err := s.cache.Character(ctx, id)
// 	if err != nil {
// 		entry.WithError(err).Error("failed to fetch character from cache")
// 		return nil, fmt.Errorf("failed to fetch character from cache")
// 	}

// 	if character == nil {
// 		cached = false
// 		character, err = s.character.Character(ctx, id)
// 		if err != nil && !errors.Is(err, sql.ErrNoRows) {
// 			entry.WithError(err).Error("failed to fetch character from DB")
// 			return nil, fmt.Errorf("failed to fetch character from DB")
// 		}

// 		if character == nil || errors.Is(err, sql.ErrNoRows) {
// 			exists = false
// 			character = &athena.Character{ID: id}
// 			err = s.esi.ResetEtag(ctx, etag)
// 			if err != nil {
// 				entry.WithError(err).Error("failed to reset etag")
// 				return nil, fmt.Errorf("failed to reset etag")
// 			}
// 		}

// 	}

// 	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

// 		if !cached {
// 			err = s.cache.SetCharacter(ctx, character)
// 			if err != nil {
// 				entry.WithError(err).Error("failed to cache character")
// 			}
// 		}

// 		return character, nil

// 	}

// 	character, _, _, err = s.esi.GetCharacter(ctx, character)
// 	if err != nil {
// 		entry.WithError(err).Error("Failed to fetch character from ESI")
// 		return nil, fmt.Errorf("Failed to fetch character from ESI")
// 	}

// 	// _, _ = s.CharacterCorporationHistory(ctx, character)

// 	switch exists {
// 	case true:
// 		character, err = s.character.UpdateCharacter(ctx, character.ID, character)
// 		if err != nil {
// 			entry.WithError(err).Error("Failed to update character in database")
// 			return nil, err
// 		}
// 	case false:
// 		character, err = s.character.CreateCharacter(ctx, character)
// 		if err != nil {
// 			entry.WithError(err).Error("Failed to create character in database")
// 			return nil, err
// 		}

// 	}

// 	err = s.cache.SetCharacter(ctx, character)
// 	if err != nil {
// 		entry.WithError(err).Error("failed to cache character")
// 	}

// 	return character, err

// }

// func (s *service) Characters(ctx context.Context, operators []*athena.Operator) ([]*athena.Character, error) {

// 	characters, err := s.character.Characters(ctx, operators...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return characters, nil

// }

// func (s *service) CharacterCorporationHistory(ctx context.Context, operators []*athena.Operator) ([]*athena.CharacterCorporationHistory, error) {

// 	// entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
// 	// 	"service": serviceIdentifier,
// 	// 	"method":  "CharacterCorporationHistory",
// 	// })

// 	// etag, err := s.esi.Etag(ctx, esi.GetCharacterCorporationHistory, esi.ModWithCharacter(character))
// 	// if err != nil {
// 	// 	entry.WithError(err).Error("failed to fetch etag object")
// 	// 	return nil, fmt.Errorf("failed to fetch etag object")
// 	// }

// 	// exists := true
// 	// cached := true

// 	// history, err := s.cache.CharacterCorporationHistory(ctx, character.ID)
// 	// if err != nil {
// 	// 	entry.WithError(err).Error("failed to cache character")
// 	// 	return nil, fmt.Errorf("failed to cache character")
// 	// }

// 	// if history == nil {
// 	// 	cached = false
// 	// 	history, err = s.character.CharacterCorporationHistory(ctx, athena.NewEqualOperator("character_id", character.ID))
// 	// 	if err != nil && !errors.Is(err, sql.ErrNoRows) {
// 	// 		entry.WithError(err).Error("failed to fetch character from DB")
// 	// 		return nil, fmt.Errorf("failed to fetch character from DB")
// 	// 	}

// 	// 	if history == nil || errors.Is(err, sql.ErrNoRows) {
// 	// 		exists = false
// 	// 		history = make([]*athena.CharacterCorporationHistory, 0)
// 	// 	}
// 	// }

// 	// if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

// 	// 	if !cached {
// 	// 		err = s.cache.SetCharacterCorporationHistory(ctx, character.ID, history)
// 	// 		if err != nil {
// 	// 			entry.WithError(err).Error("failed to cache corporation history")
// 	// 		}
// 	// 	}

// 	// 	return history, nil
// 	// }

// 	// newHistory, _, _, err := s.esi.GetCharacterCorporationHistory(ctx, character, make([]*athena.CharacterCorporationHistory, 0))
// 	// if err != nil {
// 	// 	entry.WithError(err).Error("Failed to fetch corporation history from ESI")
// 	// 	return nil, fmt.Errorf("Failed to fetch corporation history from ESI")
// 	// }

// 	// if len(newHistory) > 0 {
// 	// 	s.resolveHistoryAttributes(ctx, newHistory)
// 	// 	history, err := s.diffAndUpdateHistory(ctx, character, history, newHistory)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to diff and update corporation history")
// 	// 	}

// 	// 	err = s.cache.SetCharacterCorporationHistory(ctx, character.ID, history)
// 	// 	if err != nil {
// 	// 		entry.WithError(err).Error("failed to cache corporation history")
// 	// 	}
// 	// }

// 	return nil, nil
// }

// func (s *service) resolveHistoryAttributes(ctx context.Context, history []*athena.CharacterCorporationHistory) {

// 	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
// 		"service": serviceIdentifier,
// 		"method":  "resolveHistoryAttributes",
// 	})

// 	for _, record := range history {
// 		_, err := s.corporation.Corporation(ctx, record.CorporationID)
// 		if err != nil {
// 			entry.WithError(err).WithFields(logrus.Fields{
// 				"record_id":      record.RecordID,
// 				"corporation_id": record.CorporationID,
// 			}).Error("failed to resolve corporation record in character history")
// 		}
// 	}

// }

// func (s *service) diffAndUpdateHistory(ctx context.Context, character *athena.Character, old []*athena.CharacterCorporationHistory, new []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, error) {

// 	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
// 		"character_id": character.ID,
// 		"service":      serviceIdentifier,
// 		"method":       "diffAndUpdateHistory",
// 	})

// 	recordsToCreate := make([]*athena.CharacterCorporationHistory, 0)

// 	oldRecordMap := make(map[uint]*athena.CharacterCorporationHistory)
// 	for _, record := range old {
// 		oldRecordMap[record.RecordID] = record
// 	}

// 	for _, record := range new {
// 		if _, ok := oldRecordMap[record.RecordID]; !ok {
// 			recordsToCreate = append(recordsToCreate, record)
// 		}
// 	}

// 	var final = make([]*athena.CharacterCorporationHistory, 0)
// 	if len(recordsToCreate) > 0 {
// 		createdRecords, err := s.character.CreateCharacterCorporationHistory(ctx, character.ID, recordsToCreate)
// 		if err != nil {
// 			entry.WithError(err).Error("failed to create character corporation history records in db")
// 		}

// 		final = append(final, createdRecords...)
// 	}

// 	return final, nil

// }
