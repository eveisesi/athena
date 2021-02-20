package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type locationInterface interface {
	GetCharacterLocation(ctx context.Context, characterID uint, token string) (*athena.MemberLocation, *athena.Etag, *http.Response, error)
	GetCharacterOnline(ctx context.Context, characterID uint, token string) (*athena.MemberOnline, *athena.Etag, *http.Response, error)
	GetCharacterShip(ctx context.Context, characterID uint, token string) (*athena.MemberShip, *athena.Etag, *http.Response, error)
}

// GetCharacterLocation makes an HTTP GET Request to the /characters/{character_id}/location endpoint for
// information about the provided members current location
//
// Documentation: https://esi.evetech.net/ui/#/Location/get_characters_character_id_location
// Version: v1
// Cache: 5 secs
func (s *service) GetCharacterLocation(ctx context.Context, characterID uint, token string) (*athena.MemberLocation, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterLocation]

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

	var location = new(athena.MemberLocation)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, location)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return location, etag, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return location, etag, res, nil

}

func characterLocationsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterLocation.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterLocationsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterLocation].Path, mods.characterID)

}

func (s *service) GetCharacterOnline(ctx context.Context, characterID uint, token string) (*athena.MemberOnline, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterOnline]

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

	var online = new(athena.MemberOnline)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, online)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return online, etag, res, fmt.Errorf("failed to fetch online for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return online, etag, res, nil

}

func characterOnlinesKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterOnline.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterOnlinesPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterOnline].Path, mods.characterID)

}

func (s *service) GetCharacterShip(ctx context.Context, characterID uint, token string) (*athena.MemberShip, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterShip]

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

	var ship = new(athena.MemberShip)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, ship)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return ship, etag, res, fmt.Errorf("failed to fetch ship for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return ship, etag, res, nil

}

func characterShipsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterShip.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterShipsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterShip].Path, mods.characterID)

}
