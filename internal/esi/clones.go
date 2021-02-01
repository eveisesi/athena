package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/athena"
	"github.com/volatiletech/null"
)

type CharacterClones struct {
	HomeLocation struct {
		LocationID   uint64 `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"home_location"`
	JumpClones []struct {
		Implants     []uint      `json:"implants"`
		JumpCloneID  uint        `json:"jump_clone_id"`
		LocationID   uint64      `json:"location_id"`
		LocationType string      `json:"location_type"`
		name         null.String `json:"name"`
	} `json:"jump_clones"`
	LastCloneJumpDate     null.Time `json:"last_clone_jump_date"`
	LastStationChangeDate null.Time `json:"last_station_change_date"`
}

// GetCharacterClones makes an HTTP GET Request to the /characters/{character_id}/clones/ endpoint for
// information about the provided members clones
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_clones
// Version: v3
// Cache: 120 (2 min)
func (s *service) GetCharacterClones(
	ctx context.Context,
	member *athena.Member,
	clone *athena.MemberHomeClone,
	clones []*athena.MemberJumpClone,
) (
	*athena.MemberHomeClone,
	[]*athena.MemberJumpClone,
	*http.Response,
	error,
) {

	endpoint := s.endpoints[GetCharacterClones.Name]

	mods := s.modifiers(ModWithMember(member))

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
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	characterClones := new(CharacterClones)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, characterClones)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		spew.Dump(characterClones)

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return clone, clones, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return clone, clones, res, nil

}

func (s *service) characterClonesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterClones.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterClonesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterClones.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterClonesEndpoint() *endpoint {

	GetCharacterClones.KeyFunc = s.characterClonesKeyFunc
	GetCharacterClones.PathFunc = s.characterClonesPathFunc
	return GetCharacterClones

}

// GetCharacterImplants makes an HTTP GET Request to the /characters/{character_id}/implants/ endpoint for
// information about the provided members implants
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_implants
// Version: v1
// Cache: 120 (2 min)
func (s *service) GetCharacterImplants(ctx context.Context, member *athena.Member, ids []uint) ([]uint, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterImplants.Name]

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
		err = json.Unmarshal(b, ids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return ids, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return ids, res, nil

}

func (s *service) characterImplantsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterImplants.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterImplantsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterImplants.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterImplantsEndpoint() *endpoint {

	GetCharacterImplants.KeyFunc = s.characterImplantsKeyFunc
	GetCharacterImplants.PathFunc = s.characterImplantsPathFunc
	return GetCharacterImplants

}
