package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type clonesInterface interface {
	GetCharacterClones(ctx context.Context, characterID uint, token string) (*athena.MemberClones, *athena.Etag, *http.Response, error)
	GetCharacterImplants(ctx context.Context, characterID uint, token string) ([]uint, *athena.Etag, *http.Response, error)
}

// GetCharacterClones makes an HTTP GET Request to the /characters/{character_id}/clones/ endpoint for
// information about the provided members clones
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_clones
// Version: v3
// Cache: 120 (2 min)
func (s *service) GetCharacterClones(ctx context.Context, characterID uint, token string) (*athena.MemberClones, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterClones]

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

	var clones = new(athena.MemberClones)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, clones)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		clones.MemberID = characterID

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return clones, etag, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)

	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return clones, etag, res, nil

}

func characterClonesKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterClones.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterClonesPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterClones].Path, mods.characterID)

}

// GetCharacterImplants makes an HTTP GET Request to the /characters/{character_id}/implants/ endpoint for
// information about the provided members implants
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_implants
// Version: v1
// Cache: 120 (2 min)
func (s *service) GetCharacterImplants(ctx context.Context, characterID uint, token string) ([]uint, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterImplants]

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

	var ids = make([]uint, 0, 11)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &ids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return ids, etag, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return ids, etag, res, nil

}

func characterImplantsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterImplants.String(), strconv.FormatUint(uint64(mods.characterID), 10))
}

func characterImplantsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterImplants].Path, mods.characterID)

}
