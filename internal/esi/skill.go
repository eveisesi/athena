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

func (s *service) GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberSkillAttributes) (*athena.MemberSkillAttributes, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterAttributes.Name]

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

func (s *service) characterAttributesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterAttributes.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterAttributesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterAttributes.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterAttributesEndpoint() *endpoint {
	GetCharacterAttributes.KeyFunc = s.characterAttributesKeyFunc
	GetCharacterAttributes.PathFunc = s.characterAttributesPathFunc
	return GetCharacterAttributes
}

func (s *service) GetCharacterSkills(ctx context.Context, member *athena.Member, meta *athena.MemberSkillMeta) (*athena.MemberSkillMeta, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterSkills.Name]

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
		err = json.Unmarshal(b, &meta)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return meta, res, fmt.Errorf("failed to fetch skills for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return meta, res, nil

}

func (s *service) characterSkillsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkills.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterSkillsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterSkills.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterSkillsEndpoint() *endpoint {
	GetCharacterSkills.KeyFunc = s.characterSkillsKeyFunc
	GetCharacterSkills.PathFunc = s.characterSkillsPathFunc
	return GetCharacterSkills
}

func (s *service) GetCharacterSkillQueue(ctx context.Context, member *athena.Member, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterSkillQueue.Name]

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
		err = json.Unmarshal(b, &queue)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return queue, res, fmt.Errorf("failed to fetch skill queue for character %d, received status code of %d", member.CharacterID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return queue, res, nil

}

func (s *service) characterSkillQueueKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkillQueue.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterSkillQueuePathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterSkillQueue.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterSkillQueueEndpoint() *endpoint {
	GetCharacterSkillQueue.KeyFunc = s.characterSkillQueueKeyFunc
	GetCharacterSkillQueue.PathFunc = s.characterSkillQueuePathFunc
	return GetCharacterSkillQueue
}
