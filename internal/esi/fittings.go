package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

func (s *service) GetCharacterFittings(ctx context.Context, characterID uint, token string) ([]*athena.MemberFitting, *http.Response, error) {

	endpoint := endpoints[GetCharacterFittings]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	var fittings = make([]*athena.MemberFitting, 0, 250)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &fittings)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return fittings, res, fmt.Errorf("failed to fetch fittings for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return fittings, res, nil

}

func characterFittingsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterFittings.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterFittingsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterFittings].Path, mods.characterID)

}
