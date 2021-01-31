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

// GetCharacterClones makes an HTTP GET Request to the /characters/{character_id}/clones/ endpoint for
// information about the provided members clones
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_clones
// Version: v3
// Cache: 120 (2 min)
func (s *service) GetCharacterClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterClones.Name]

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

func (s *service) characterClonesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterClones.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterClonesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterClones.FmtPath, mods.member.CharacterID),
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
func (s *service) GetCharacterImplants(ctx context.Context, member *athena.Member, implants *athena.MemberImplants) (*athena.MemberImplants, *http.Response, error) {

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

func (s *service) characterImplantsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterImplants.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterImplantsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterImplants.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterImplantsEndpoint() *endpoint {

	GetCharacterImplants.KeyFunc = s.characterImplantsKeyFunc
	GetCharacterImplants.PathFunc = s.characterImplantsPathFunc
	return GetCharacterImplants

}
