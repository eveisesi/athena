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
	HeadCharacterContacts(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error)
	GetCharacterContacts(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberContact, *athena.Etag, *http.Response, error)
	GetCharacterContactLabels(ctx context.Context, characterID uint, token string) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterContacts(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContacts]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithPage(page))

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
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag: %w", err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterContacts(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberContact, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContacts]

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
		WithPage(page),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, etag, res, nil
	}

	contacts := make([]*athena.MemberContact, 0, 1024) // ESI Specification states a max of 1024 items can be returned per page
	err = json.Unmarshal(b, &contacts)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return contacts, nil, res, nil

}

func characterContactsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	params := make([]string, 0, 3)
	params = append(params, GetCharacterContacts.String(), strconv.FormatUint(uint64(mods.characterID), 10))

	if mods.page > 0 {
		params = append(params, strconv.FormatUint(uint64(mods.page), 10))
	}

	return buildKey(params...)

}

func characterContactsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterContacts].Path, mods.characterID)

}

func (s *service) GetCharacterContactLabels(ctx context.Context, characterID uint, token string) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContactLabels]

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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch contact labels for character %d, received status code of %d", characterID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, etag, res, nil
	}

	var labels = make([]*athena.MemberContactLabel, 0, 64)
	err = json.Unmarshal(b, &labels)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return labels, etag, res, nil

}

func characterContactLabelsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(
		GetCharacterContactLabels.String(),
		strconv.FormatUint(uint64(mods.characterID), 10),
	)

}

func characterContactLabelsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterContactLabels].Path, mods.characterID)

}
