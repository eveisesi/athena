package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func isCharacterValid(r *athena.Character) bool {
	if r.Name == "" || r.CorporationID == 0 {
		return false
	}
	return true
}

// GetCharactersCharacterID makes a HTTP GET Request to the /characters/{character_id} endpoint
// for information about the provided character
//
// Documentation: https://esi.evetech.net/ui/#/Character/get_characters_character_id
// Version: v4
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharactersCharacterID(ctx context.Context, character *athena.Character) (*athena.Character, *http.Response, error) {

	path := s.endpoints[EndpointGetCharactersCharacterID](character)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path), WithEtag(character.Etag))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, character)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			character.Etag = etag
		}

		if !isCharacterValid(character) {
			return nil, nil, fmt.Errorf("invalid character return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return character, res, fmt.Errorf("failed to fetch character %d, received status code of %d", character.CharacterID, sc)
	}

	character.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return character, res, nil
}

func (s *service) resolveGetCharactersCharacterIDEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Character, received nil")
	}

	var thing *athena.Character
	var ok bool

	if thing, ok = obj.(*athena.Character); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Character, got %T", obj))
	}

	return fmt.Sprintf("/v4/characters/%d/", thing.CharacterID)

}
