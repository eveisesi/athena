package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type contactsInterface interface {
	GetCharacterContacts(ctx context.Context, member *athena.Member, contacts []*athena.MemberContact) ([]*athena.MemberContact, *athena.Etag, *http.Response, error)
	GetCharacterContactLabels(ctx context.Context, member *athena.Member, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterContacts(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContacts]

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.ID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterContacts(ctx context.Context, member *athena.Member, page uint) ([]*athena.MemberContact, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContacts]

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
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}
	contacts := make([]*athena.MemberContact, 0, 1024) // ESI Specification states a max of 1020 items can be returned per page

	if res.StatusCode >= http.StatusBadRequest {
		return contacts, nil, res, fmt.Errorf("failed to exec contacts head request for character %d, received status code of %d", member.ID, res.StatusCode)
	}

	err = json.Unmarshal(b, &contacts)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	etag.Etag = s.retrieveEtagHeader(res.Header)
	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return contacts, nil, res, nil

}

func characterContactsKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requirePage(mods)

	return buildKey(
		GetCharacterContacts.String(),
		strconv.FormatUint(uint64(mods.member.ID), 10),
		strconv.FormatUint(uint64(mods.page), 10),
	)

}

func characterContactsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterContacts].Path, mods.member.ID)

}

func (s *service) GetCharacterContactLabels(ctx context.Context, member *athena.Member, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContactLabels]

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

func characterContactLabelsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(
		GetCharacterContactLabels.String(),
		strconv.FormatUint(uint64(mods.member.ID), 10),
	)

}

func characterContactLabelsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterContactLabels].Path, mods.member.ID)

}
