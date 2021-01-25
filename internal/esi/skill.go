package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func (s *service) GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberSkillAttributes) (*athena.MemberSkillAttributes, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterAttributes](member)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(attributes.Etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, attributes)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			attributes.Etag = etag
		}

	case sc >= http.StatusBadRequest:
		return attributes, res, fmt.Errorf("failed to fetch attributes for character %d, received status code of %d", member.CharacterID, sc)
	}

	attributes.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return attributes, res, nil

}

func (s *service) resolveGetCharacterAttributes(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/attributes/", thing.CharacterID)

}

func (s *service) GetCharacterSkills(ctx context.Context, member *athena.Member, etag *athena.Etag, meta *athena.MemberSkillMeta) (*athena.MemberSkillMeta, *athena.Etag, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterSkills](member)

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
		err = json.Unmarshal(b, &meta)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return meta, etag, res, fmt.Errorf("failed to fetch skills for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return meta, etag, res, nil

}

func (s *service) resolveGetCharacterSkills(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/skills/", thing.CharacterID)

}

func (s *service) GetCharacterSkillQueue(ctx context.Context, member *athena.Member, etag *athena.Etag, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error) {

	path := s.endpoints[EndpointGetCharacterSkillQueue](member)

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
		err = json.Unmarshal(b, &queue)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return queue, etag, res, fmt.Errorf("failed to fetch skill queue for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return queue, etag, res, nil

}

func (s *service) resolveGetCharacterSkillQueue(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Member, received nil")
	}

	var thing *athena.Member
	var ok bool

	if thing, ok = obj.(*athena.Member); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Member, got %T", obj))
	}

	return fmt.Sprintf("/v1/characters/%d/skillqueue/", thing.CharacterID)

}
