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

type clonesInterface interface {
	GetCharacterClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *athena.Etag, *http.Response, error)
	GetCharacterImplants(ctx context.Context, member *athena.Member, ids []uint) ([]uint, *athena.Etag, *http.Response, error)
}

// GetCharacterClones makes an HTTP GET Request to the /characters/{character_id}/clones/ endpoint for
// information about the provided members clones
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_clones
// Version: v3
// Cache: 120 (2 min)
func (s *service) GetCharacterClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterClones]

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

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, clones)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return clones, etag, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return clones, etag, res, nil

}

func characterClonesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterClones.String(), strconv.Itoa(int(mods.member.ID)))
}

func characterClonesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterClones].Path, mods.member.ID),
	}

	return u.String()

}

// GetCharacterImplants makes an HTTP GET Request to the /characters/{character_id}/implants/ endpoint for
// information about the provided members implants
//
// Documentation: https://esi.evetech.net/ui/#/Clones/get_characters_character_id_implants
// Version: v1
// Cache: 120 (2 min)
func (s *service) GetCharacterImplants(ctx context.Context, member *athena.Member, ids []uint) ([]uint, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterImplants]

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

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &ids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return ids, etag, res, fmt.Errorf("failed to fetch clones for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return ids, etag, res, nil

}

func characterImplantsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterImplants.String(), strconv.Itoa(int(mods.member.ID)))
}

func characterImplantsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterImplants].Path, mods.member.ID),
	}

	return u.String()

}
