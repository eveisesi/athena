package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eveisesi/athena"
)

func isCharacterValid(r *athena.Character) bool {
	if r.Name == "" || r.CorporationID == 0 {
		return false
	}
	return true
}

// GetCharacter makes a HTTP GET Request to the /characters/{character_id} endpoint
// for information about the provided character
//
// Documentation: https://esi.evetech.net/ui/#/Character/get_characters_character_id
// Version: v4
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharacter(ctx context.Context, character *athena.Character) (*athena.Character, *http.Response, error) {

	endpoint := s.endpoints[GetCharacter.Name]

	mods := s.modifiers(ModWithCharacter(character))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path), WithEtag(etag))
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

func (s *service) newGetCharacterEndpoint() *endpoint {
	GetCharacter.KeyFunc = s.characterKeyFunc
	GetCharacter.PathFunc = s.characterPathFunc
	return GetCharacter
}

func (s *service) characterKeyFunc(mods *modifiers) string {

	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil for character instead")
	}

	return buildKey(GetCharacter.Name, strconv.Itoa(int(mods.character.CharacterID)))
}

func (s *service) characterPathFunc(mods *modifiers) string {

	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil for character instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacter.FmtPath, mods.character.CharacterID),
	}

	return u.String()

}
