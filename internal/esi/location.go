package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

// GetCharacterLocation makes an HTTP GET Request to the /characters/{character_id}/location endpoint for
// information about the provided members current location
//
// Documentation: https://esi.evetech.net/ui/#/Location/get_characters_character_id_location
// Version: v1
// Cache: 5 secs
func (s *service) GetCharacterLocation(ctx context.Context, member *athena.Member, location *athena.MemberLocation) (*athena.MemberLocation, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterLocation](member)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(location.Etag),
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

func (s *service) resolveGetCharacterLocationEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/location/", thing.CharacterID)

}

func (s *service) GetCharacterOnline(ctx context.Context, member *athena.Member, online *athena.MemberOnline) (*athena.MemberOnline, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterOnline](member)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(online.Etag),
		WithAuthorization(member.AccessToken),
	)
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

func (s *service) resolveGetCharacterOnlineEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v2/characters/%d/online/", thing.CharacterID)

}

func (s *service) GetCharacterShip(ctx context.Context, member *athena.Member, ship *athena.MemberShip) (*athena.MemberShip, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterShip](member)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(ship.Etag),
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

func (s *service) resolveGetCharacterShipEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/ship/", thing.CharacterID)

}
