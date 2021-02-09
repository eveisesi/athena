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

type characterInterface interface {
	GetCharacter(ctx context.Context, character *athena.Character) (*athena.Character, *athena.Etag, *http.Response, error)
	GetCharacterCorporationHistory(ctx context.Context, character *athena.Character, history []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, *athena.Etag, *http.Response, error)
}

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
func (s *service) GetCharacter(ctx context.Context, character *athena.Character) (*athena.Character, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacter.Name]

	mods := s.modifiers(ModWithCharacter(character))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, character)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

		if !isCharacterValid(character) {
			return nil, nil, nil, fmt.Errorf("invalid character return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return character, etag, res, fmt.Errorf("failed to fetch character %d, received status code of %d", character.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return character, etag, res, nil

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

	return buildKey(GetCharacter.Name, strconv.Itoa(int(mods.character.ID)))
}

func (s *service) characterPathFunc(mods *modifiers) string {

	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil for character instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacter.FmtPath, mods.character.ID),
	}

	return u.String()

}

// GetCharacterCorporationHistory makes a HTTP GET Request to the /v1/characters/{character_id}/corporationhistory/ endpoint
// for information about the provided characters corporation history
//
// Documentation: https://esi.evetech.net/ui/?version=_latest#/Character/get_characters_character_id_corporationhistory
// Version: v1
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharacterCorporationHistory(ctx context.Context, character *athena.Character, history []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterCorporationHistory.Name]

	mods := s.modifiers(ModWithCharacter(character))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &history)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return history, etag, res, fmt.Errorf("failed to fetch character history %d, received status code of %d", character.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return history, etag, res, nil

}

func (s *service) newGetCharacterCorporationHistoryEndpoint() *endpoint {

	GetCharacterCorporationHistory.KeyFunc = s.characterCorporationHistoryKeyFunc
	GetCharacterCorporationHistory.PathFunc = s.characterCorporationHistoryPathFunc
	return GetCharacterCorporationHistory

}

func (s *service) characterCorporationHistoryKeyFunc(mods *modifiers) string {

	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil for character instead")
	}

	return buildKey(GetCharacterCorporationHistory.Name, strconv.Itoa(int(mods.character.ID)))

}

func (s *service) characterCorporationHistoryPathFunc(mods *modifiers) string {

	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil for character instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterCorporationHistory.FmtPath, mods.character.ID),
	}

	return u.String()

}
