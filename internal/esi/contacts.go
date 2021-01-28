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

func (s *service) GetCharacterContacts(ctx context.Context, member *athena.Member, contacts []*athena.MemberContact) ([]*athena.MemberContact, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterContacts.Name]

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
		err = json.Unmarshal(b, &contacts)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return contacts, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return contacts, res, nil

}

func (s *service) characterContactsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Alliance to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterContacts.Name, strconv.Itoa(int(mods.member.CharacterID)))
}

func (s *service) characterContactsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Alliance to be provided, received nil for alliance instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContacts.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterContactsEndpoint() *endpoint {

	GetCharacterContacts.KeyFunc = s.characterContactsKeyFunc
	GetCharacterContacts.PathFunc = s.characterContactsPathFunc
	return GetCharacterContacts

}

func (s *service) GetCharacterContactLabels(ctx context.Context, member *athena.Member, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterContactLabels.Name]

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
		err = json.Unmarshal(b, &labels)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return labels, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return labels, res, nil

}

func (s *service) characterContactLabelsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Alliance to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterContactLabels.Name, strconv.Itoa(int(mods.member.CharacterID)))
}

func (s *service) characterContactLabelsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Alliance to be provided, received nil for alliance instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContactLabels.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterContactLabelsEndpoint() *endpoint {

	GetCharacterContactLabels.KeyFunc = s.characterContactLabelsKeyFunc
	GetCharacterContactLabels.PathFunc = s.characterContactLabelsPathFunc
	return GetCharacterContactLabels

}
