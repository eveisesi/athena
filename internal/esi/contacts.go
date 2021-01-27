package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func (s *service) GetCharactersCharacterIDContacts(ctx context.Context, member *athena.Member, contacts []*athena.MemberContact) ([]*athena.MemberContact, *http.Response, error) {

	iterator := 0
	for {
		path := s.endpoints[EndpointGetCharactersCharacterIDContacts](member)

		b, res, err := s.request(
			ctx,
			WithMethod(http.MethodGet),
			WithPath(path),
			WithEtag(etag.Etag),
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
				return nil, nil, nil, err
			}

			etag.Etag = s.retrieveEtagHeader(res.Header)

		case sc >= http.StatusBadRequest:
			return contacts, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.CharacterID, sc)
		}

		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	}

	return contacts, res, nil

}

func (s *service) GetCharacterContactLabels(ctx context.Context, member *athena.Member, etag *athena.Etag, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterContactLabels](member)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
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
		return labels, etag, res, fmt.Errorf("failed to fetch location for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return labels, etag, res, nil

}

func (s *service) resolveGetCharacterContactsEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v2/characters/%d/contacts/", thing.CharacterID)

}

func (s *service) resolveGetCharacterContactLabelsEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/contacts/labels", thing.CharacterID)

}
