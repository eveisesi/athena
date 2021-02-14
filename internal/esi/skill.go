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

type skillsInterface interface {
	GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberAttributes) (*athena.MemberAttributes, *athena.Etag, *http.Response, error)
	GetCharacterSkills(ctx context.Context, member *athena.Member, meta *athena.MemberSkills) (*athena.MemberSkills, *athena.Etag, *http.Response, error)
	GetCharacterSkillQueue(ctx context.Context, member *athena.Member, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error)
}

func (s *service) GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberAttributes) (*athena.MemberAttributes, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterAttributes]

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

func characterAttributesKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterAttributes.String(), strconv.Itoa(int(mods.member.ID)))
}

func characterAttributesPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterAttributes].Path, mods.member.ID),
	}

	return u.String()

}

func (s *service) GetCharacterSkills(ctx context.Context, member *athena.Member, skills *athena.MemberSkills) (*athena.MemberSkills, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterSkills]

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

func characterSkillsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkills.String(), strconv.Itoa(int(mods.member.ID)))
}

func characterSkillsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterSkills].Path, mods.member.ID),
	}

	return u.String()

}

func (s *service) GetCharacterSkillQueue(ctx context.Context, member *athena.Member, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterSkillQueue]

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

func characterSkillQueueKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for alliance instead")
	}

	return buildKey(GetCharacterSkillQueue.String(), strconv.Itoa(int(mods.member.ID)))
}

func characterSkillQueuePathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterSkillQueue].Path, mods.member.ID),
	}

	return u.String()

}
