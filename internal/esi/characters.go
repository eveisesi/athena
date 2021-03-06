package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type characterInterface interface {
	GetCharacter(ctx context.Context, characterID uint) (*athena.Character, *athena.Etag, *http.Response, error)
	GetCharacterCorporationHistory(ctx context.Context, characterID uint) ([]*athena.CharacterCorporationHistory, *athena.Etag, *http.Response, error)
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
func (s *service) GetCharacter(ctx context.Context, characterID uint) (*athena.Character, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacter]

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
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch character %d, received status code of %d", characterID, res.StatusCode)
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

	var character = new(athena.Character)
	err = json.Unmarshal(b, character)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	character.ID = characterID

	if !isCharacterValid(character) {
		return nil, nil, nil, fmt.Errorf("invalid character return from esi, missing name or ticker")
	}

	return character, etag, res, nil

}

func characterKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacter.String(), strconv.Itoa(int(mods.characterID)))
}

func characterPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacter].Path, mods.characterID)

}

// GetCharacterCorporationHistory makes a HTTP GET Request to the /v1/characters/{character_id}/corporationhistory/ endpoint
// for information about the provided characters corporation history
//
// Documentation: https://esi.evetech.net/ui/?version=_latest#/Character/get_characters_character_id_corporationhistory
// Version: v1
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharacterCorporationHistory(ctx context.Context, characterID uint) ([]*athena.CharacterCorporationHistory, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterCorporationHistory]

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
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch corporation history for character %d, received status code of %d", characterID, res.StatusCode)
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

	var history = make([]*athena.CharacterCorporationHistory, 0, 512)
	err = json.Unmarshal(b, &history)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return history, etag, res, nil

}

func characterCorporationHistoryKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterCorporationHistory.String(), strconv.Itoa(int(mods.characterID)))

}

func characterCorporationHistoryPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterCorporationHistory].Path, mods.characterID)

}
