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

// GetCharacterLocation makes an HTTP GET Request to the /characters/{character_id}/location endpoint for
// information about the provided members current location
//
// Documentation: https://esi.evetech.net/ui/#/Location/get_characters_character_id_location
// Version: v1
// Cache: 5 secs
func (s *service) GetCharacterLocation(ctx context.Context, member *athena.Member, location *athena.MemberLocation) (*athena.MemberLocation, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterLocation.Name]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, location)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			location.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return location, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.CharacterID, sc)
	}

	location.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return location, res, nil

}

func (s *service) characterLocationsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterLocation.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterLocationsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterLocation.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterLocationEndpoint() *endpoint {

	GetCharacterLocation.KeyFunc = s.characterLocationsKeyFunc
	GetCharacterLocation.PathFunc = s.characterLocationsPathFunc
	return GetCharacterLocation

}

func (s *service) GetCharacterOnline(ctx context.Context, member *athena.Member, online *athena.MemberOnline) (*athena.MemberOnline, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterOnline.Name]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, online)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			online.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return online, res, fmt.Errorf("failed to fetch online for character %d, received status code of %d", member.CharacterID, sc)
	}

	online.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return online, res, nil

}

func (s *service) characterOnlinesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterOnline.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterOnlinesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterOnline.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterOnlineEndpoint() *endpoint {

	GetCharacterOnline.KeyFunc = s.characterOnlinesKeyFunc
	GetCharacterOnline.PathFunc = s.characterOnlinesPathFunc
	return GetCharacterOnline

}

func (s *service) GetCharacterShip(ctx context.Context, member *athena.Member, ship *athena.MemberShip) (*athena.MemberShip, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterShip.Name]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, ship)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			ship.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return ship, res, fmt.Errorf("failed to fetch ship for character %d, received status code of %d", member.CharacterID, sc)
	}

	ship.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return ship, res, nil

}

func (s *service) characterShipsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterShip.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterShipsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterShip.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterShipEndpoint() *endpoint {

	GetCharacterShip.KeyFunc = s.characterShipsKeyFunc
	GetCharacterShip.PathFunc = s.characterShipsPathFunc
	return GetCharacterShip

}
