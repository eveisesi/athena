package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type fittingInterface interface {
	HeadCharacterFittings(ctx context.Context, characterID uint, token string) (*athena.Etag, *http.Response, error)
	GetCharacterFittings(ctx context.Context, characterID uint, token string) ([]*athena.MemberFitting, *athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterFittings(ctx context.Context, characterID uint, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterFittings]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to fetch fittings for character %d, received status code of %d", characterID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag: %w", err)
		}
		return etag, res, nil
	}

	etag.Etag = RetrieveEtagHeader(res.Header)

	return etag, res, nil

}

func (s *service) GetCharacterFittings(ctx context.Context, characterID uint, token string) ([]*athena.MemberFitting, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterFittings]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch fittings for character %d, received status code of %d", characterID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, etag, res, nil
	}

	var fittings = make([]*athena.MemberFitting, 0, 250)
	err = json.Unmarshal(b, &fittings)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return fittings, etag, res, nil

}

func characterFittingsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterFittings.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterFittingsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterFittings].Path, mods.characterID)

}
