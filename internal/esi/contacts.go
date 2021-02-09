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

type contactsInterface interface {
	GetCharacterContacts(ctx context.Context, member *athena.Member, contacts []*athena.MemberContact) ([]*athena.MemberContact, *athena.Etag, *http.Response, error)
	GetCharacterContactLabels(ctx context.Context, member *athena.Member, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error)
}

func (s *service) GetCharacterContacts(ctx context.Context, member *athena.Member, contacts []*athena.MemberContact) ([]*athena.MemberContact, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterContacts.Name]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return contacts, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", member.ID, res.StatusCode)
	} else {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		if res.StatusCode == http.StatusNotModified {
			return contacts, etag, res, nil
		}
	}

	pages := s.retrieveXPagesFromHeader(res.Header)
	if pages == 0 {
		return nil, nil, nil, fmt.Errorf("received 0 for X-Pages on request %s, expected number greater than 0", path)
	}

	for i := 1; i <= pages; i++ {

		pageContacts := make([]*athena.MemberContact, 0)

		mods := s.modifiers(ModWithMember(member), ModWithPage(&i))

		pageEtag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
		if err != nil {
			return nil, nil, nil, err
		}

		path := endpoint.PathFunc(mods)

		b, res, err := s.request(
			ctx,
			WithMethod(http.MethodGet),
			WithPath(path),
			WithEtag(pageEtag),
			WithPage(i),
			WithAuthorization(member.AccessToken),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageContacts)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, nil, err
			}

			contacts = append(contacts, pageContacts...)

			pageEtag.Etag = s.retrieveEtagHeader(res.Header)

		case sc >= http.StatusBadRequest:
			return contacts, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", member.ID, sc)
		}

		pageEtag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, pageEtag.EtagID, pageEtag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
		}
	}

	return contacts, etag, res, nil

}

func (s *service) characterContactsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	param := append(make([]string, 0), GetCharacterContacts.Name, strconv.Itoa(int(mods.member.ID)))

	if mods.page != nil {
		param = append(param, strconv.Itoa(*mods.page))
	}

	return buildKey(param...)

}

func (s *service) characterContactsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContacts.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterContactsEndpoint() *endpoint {

	GetCharacterContacts.KeyFunc = s.characterContactsKeyFunc
	GetCharacterContacts.PathFunc = s.characterContactsPathFunc
	return GetCharacterContacts

}

func (s *service) GetCharacterContactLabels(ctx context.Context, member *athena.Member, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterContactLabels.Name]

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
		err = json.Unmarshal(b, &labels)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return labels, etag, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return labels, etag, res, nil

}

func (s *service) characterContactLabelsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterContactLabels.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterContactLabelsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContactLabels.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterContactLabelsEndpoint() *endpoint {

	GetCharacterContactLabels.KeyFunc = s.characterContactLabelsKeyFunc
	GetCharacterContactLabels.PathFunc = s.characterContactLabelsPathFunc
	return GetCharacterContactLabels

}
