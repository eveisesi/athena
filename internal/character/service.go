package character

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	CharacterByCharacterID(ctx context.Context, id uint64, options []OptionFunc) (*athena.Character, error)
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

func (s *service) CharacterByCharacterID(ctx context.Context, id uint64, options []OptionFunc) (*athena.Character, error) {

	characters, err := s.Characters(ctx, athena.NewOperators(athena.NewEqualOperator("character_id", id)), options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	if len(characters) == 1 {
		return characters[0], nil
	}

	character, _, err := s.esi.GetCharactersCharacterID(ctx, &athena.Character{CharacterID: id})
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	character, err = s.CreateCharacter(ctx, character, options)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return character, nil

}

func (s *service) Characters(ctx context.Context, operators []*athena.Operator, options []OptionFunc) ([]*athena.Character, error) {

	opts := s.options(options)

	if !opts.skipCache {
		characters, err := s.cache.Characters(ctx, operators)
		if err != nil {
			return nil, err
		}

		if characters != nil {
			return characters, nil
		}
	}

	characters, err := s.CharacterRepository.Characters(ctx, operators...)
	if err != nil {
		return nil, err
	}

	if opts.skipCache {
		return characters, nil
	}

	err = s.cache.SetCharacters(ctx, operators, characters, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return characters, nil

}

func (s *service) CreateCharacter(ctx context.Context, character *athena.Character, options []OptionFunc) (*athena.Character, error) {

	character, err := s.CharacterRepository.CreateCharacter(ctx, character)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	opts := s.options(options)

	if opts.skipCache {
		return character, err
	}

	err = s.cache.SetCharacter(ctx, character, cache.ExpiryMinutes(30))
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return character, nil

}
