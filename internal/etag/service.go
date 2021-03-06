package etag

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
)

type Service interface {
	athena.EtagRepository
}

type service struct {
	cache cache.Service

	athena.EtagRepository
}

func NewService(cache cache.Service, etag athena.EtagRepository) Service {
	return &service{
		cache: cache,

		EtagRepository: etag,
	}
}

func (s *service) Etag(ctx context.Context, etagID string) (*athena.Etag, error) {

	etag, err := s.cache.Etag(ctx, etagID)
	if err != nil {
		return nil, err
	}

	if etag == nil {
		etag, err = s.EtagRepository.Etag(ctx, etagID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if etag == nil || errors.Is(err, sql.ErrNoRows) {
			etag = &athena.Etag{
				EtagID: etagID,
			}
			err = nil
		}
	}

	return etag, err

}

func (s *service) UpdateEtag(ctx context.Context, etagID string, etag *athena.Etag) (*athena.Etag, error) {
	var err error

	etag, err = s.EtagRepository.InsertEtag(ctx, etag)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetEtag(
		ctx, etag.EtagID, etag,
		time.Since(
			etag.CachedUntil,
		),
	)

	return etag, nil

}
