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

func (s *service) GetCharacterFittings(ctx context.Context, member *athena.Member, fittings []*athena.MemberFitting) ([]*athena.MemberFitting, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterFittings.Name]

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
		err = json.Unmarshal(b, &fittings)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return fittings, res, fmt.Errorf("failed to fetch fittings for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
	}

	return fittings, res, nil

}

func (s *service) newGetCharacterFittingsEndpoint() *endpoint {

	GetCharacterFittings.KeyFunc = func(mods *modifiers) string {

		if mods.member == nil {
			panic("expected type *athena.Member to be provided, received nil for member instead")
		}

		param := append(make([]string, 0), GetCharacterFittings.Name, strconv.Itoa(int(mods.member.ID)))

		if mods.page != nil {
			param = append(param, strconv.Itoa(*mods.page))
		}

		return buildKey(param...)
	}

	GetCharacterFittings.PathFunc = func(mods *modifiers) string {

		if mods.member == nil {
			panic("expected type *athena.Member to be provided, received nil for member instead")
		}

		u := url.URL{
			Path: fmt.Sprintf(GetCharacterFittings.FmtPath, mods.member.ID),
		}

		return u.String()

	}

	return GetCharacterFittings

}
