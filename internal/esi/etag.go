package esi

import (
	"context"

	"github.com/eveisesi/athena"
)

func (s *service) Etag(ctx context.Context, endpoint *endpoint, modifierFunc ...modifierFunc) (*athena.Etag, error) {

	mods := s.modifiers(modifierFunc...)

	if mods.page != nil {
		mods.page = nil
	}

	// Use endpoint to resolve endpoint resolver function
	key := endpoint.KeyFunc(mods)

	etags, err := s.etag.Etags(ctx, athena.NewLikeOperator("etagID", key))
	if err != nil {
		return nil, err
	}

	var oldest *athena.Etag
	for _, etag := range etags {
		if oldest == nil {
			oldest = etag
			continue
		} else if etag.CachedUntil.After(oldest.CachedUntil) {
			oldest = etag
		}
	}

	return oldest, nil
}
