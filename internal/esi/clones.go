package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

// GetCharactersCharacterIDClones makes an HTTP GET Request to the /characters/{character_id}/clones/ endpoint for
// information about the provided members clones
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_clones
// Version: v3
// Cache: 120 (2 min)
func (s *service) GetCharactersCharacterIDClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *http.Response, error) {

	path := fmt.Sprintf("/v3/characters/%d/clones/", member.CharacterID)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(clones.Etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, clones)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			clones.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return clones, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.CharacterID, sc)
	}

	clones.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return clones, res, nil

}

// GetCharactersCharacterIDImplants makes an HTTP GET Request to the /characters/{character_id}/implants/ endpoint for
// information about the provided members implants
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_implants
// Version: v1
// Cache: 120 (2 min)
func (s *service) GetCharactersCharacterIDImplants(ctx context.Context, member *athena.Member, implants *athena.MemberImplants) (*athena.MemberImplants, *http.Response, error) {

	path := fmt.Sprintf("/v1/characters/%d/implants/", member.CharacterID)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(implants.Etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		implants.Raw = make([]int, 0)
		err = json.Unmarshal(b, &implants.Raw)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			implants.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return implants, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.CharacterID, sc)
	}

	implants.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return implants, res, nil

}
