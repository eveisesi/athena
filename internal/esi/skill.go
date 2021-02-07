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

type skillInterface interface {
	GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberAttributes) (*athena.MemberAttributes, *athena.Etag, *http.Response, error)
	GetCharacterSkills(ctx context.Context, member *athena.Member, meta *athena.MemberSkills) (*athena.MemberSkills, *athena.Etag, *http.Response, error)
	GetCharacterSkillQueue(ctx context.Context, member *athena.Member, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error)
}

func (s *service) GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberAttributes) (*athena.MemberAttributes, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterAttributes.Name]

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
		err = json.Unmarshal(b, attributes)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return attributes, etag, res, fmt.Errorf("failed to fetch attributes for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return attributes, etag, res, nil

}

func (s *service) characterAttributesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterAttributes.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterAttributesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterAttributes.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterAttributesEndpoint() *endpoint {
	GetCharacterAttributes.KeyFunc = s.characterAttributesKeyFunc
	GetCharacterAttributes.PathFunc = s.characterAttributesPathFunc
	return GetCharacterAttributes
}

func (s *service) GetCharacterSkills(ctx context.Context, member *athena.Member, skills *athena.MemberSkills) (*athena.MemberSkills, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterSkills.Name]

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
		err = json.Unmarshal(b, &skills)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return skills, etag, res, fmt.Errorf("failed to fetch skills for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return skills, etag, res, nil

}

func (s *service) characterSkillsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkills.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterSkillsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterSkills.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterSkillsEndpoint() *endpoint {
	GetCharacterSkills.KeyFunc = s.characterSkillsKeyFunc
	GetCharacterSkills.PathFunc = s.characterSkillsPathFunc
	return GetCharacterSkills
}

func (s *service) GetCharacterSkillQueue(ctx context.Context, member *athena.Member, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterSkillQueue.Name]

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
		err = json.Unmarshal(b, &queue)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return queue, etag, res, fmt.Errorf("failed to fetch skill queue for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return queue, etag, res, nil

}

func (s *service) characterSkillQueueKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkillQueue.Name, strconv.Itoa(int(mods.member.ID)))
}

func (s *service) characterSkillQueuePathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterSkillQueue.FmtPath, mods.member.ID),
	}

	return u.String()

}

func (s *service) newGetCharacterSkillQueueEndpoint() *endpoint {
	GetCharacterSkillQueue.KeyFunc = s.characterSkillQueueKeyFunc
	GetCharacterSkillQueue.PathFunc = s.characterSkillQueuePathFunc
	return GetCharacterSkillQueue
}
