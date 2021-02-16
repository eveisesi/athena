package esi

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
)

type etagInterface interface {
	Etag(ctx context.Context, endpoint endpointID, modifierFunc ...modifierFunc) (*athena.Etag, error)
	ResetEtag(ctx context.Context, etag *athena.Etag) error
}

func (s *service) Etag(ctx context.Context, endpoint endpointID, modifierFunc ...modifierFunc) (*athena.Etag, error) {

	mods := s.modifiers(modifierFunc...)
	mods.page = 0

	e := endpoints[endpoint]

	// Use endpoint to resolve endpoint resolver function
	key := e.KeyFunc(mods)

	etags, err := s.etag.Etags(ctx, athena.NewLikeOperator("etag_id", key))
	if err != nil {
		return nil, err
	}

	var oldest *athena.Etag
	for _, etag := range etags {
		if oldest == nil {
			oldest = etag
			continue
		} else if etag.CachedUntil.Before(oldest.CachedUntil) {
			oldest = etag
		}
	}

	return oldest, nil

}

func (s *service) ResetEtag(ctx context.Context, etag *athena.Etag) error {

	if etag == nil {
		return nil
	}

	_, err := s.etag.DeleteEtag(ctx, etag.EtagID)
	if err != nil {
		return fmt.Errorf("failed to reset Etag: %w", err)
	}

	err = s.cache.DeleteEtag(ctx, etag.EtagID)
	if err != nil {
		return fmt.Errorf("failed to remove Etag from Cache: %w", err)
	}

	return nil

}
