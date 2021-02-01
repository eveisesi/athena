package character

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	Character(ctx context.Context, id uint, options []OptionFunc) (*athena.Character, error)
	Characters(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Character, error)
	CreateCharacter(ctx context.Context, character *athena.Character, options []OptionFunc) (*athena.Character, error)
}

type service struct {
	cache cache.Service
	esi   esi.Service
	athena.CharacterRepository
}

func NewService(cache cache.Service, esi esi.Service, character athena.CharacterRepository) Service {
	return &service{
		cache: cache,
		esi:   esi,

		CharacterRepository: character,
	}
}

func (s *service) Character(ctx context.Context, id uint, options []OptionFunc) (*athena.Character, error) {

	character, err := s.CharacterRepository.Character(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("[Character Service] Failed to fetch character %d from db: %w", id, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	character, _, err = s.esi.GetCharacter(ctx, &athena.Character{ID: id})
	if err != nil {
		err = fmt.Errorf("[Character Service] Failed to fetch character %d from ESI: %w", id, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	character, err = s.CreateCharacter(ctx, character, options)
	if err != nil {
		err = fmt.Errorf("[Character Service] Failed to create character %d in database: %w", id, err)
	}

	return character, err

}

func (s *service) Characters(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Character, error) {

	characters, err := s.CharacterRepository.Characters(ctx, operators...)
	if err != nil {
		return nil, err
	}

	return characters, nil

}

func (s *service) CreateCharacter(ctx context.Context, character *athena.Character, options []OptionFunc) (*athena.Character, error) {
	return s.CharacterRepository.CreateCharacter(ctx, character)
}
