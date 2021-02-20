package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type skillsInterface interface {
	GetCharacterAttributes(ctx context.Context, characterID uint, token string) (*athena.MemberAttributes, *athena.Etag, *http.Response, error)
	GetCharacterSkills(ctx context.Context, characterID uint, token string) (*athena.MemberSkills, *athena.Etag, *http.Response, error)
	GetCharacterSkillQueue(ctx context.Context, characterID uint, token string) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error)
}

func (s *service) GetCharacterAttributes(ctx context.Context, characterID uint, token string) (*athena.MemberAttributes, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterAttributes]

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

	var attributes = new(athena.MemberAttributes)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, attributes)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return attributes, etag, res, fmt.Errorf("failed to fetch attributes for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return attributes, etag, res, nil

}

func characterAttributesKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterAttributes.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterAttributesPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterAttributes].Path, mods.characterID)

}

func (s *service) GetCharacterSkills(ctx context.Context, characterID uint, token string) (*athena.MemberSkills, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterSkills]

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

	var skills = new(athena.MemberSkills)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, skills)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return skills, etag, res, fmt.Errorf("failed to fetch skills for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return skills, etag, res, nil

}

func characterSkillsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterSkills.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterSkillsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterSkills].Path, mods.characterID)

}

func (s *service) GetCharacterSkillQueue(ctx context.Context, characterID uint, token string) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterSkillQueue]

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

	var queue = make([]*athena.MemberSkillQueue, 0, 51)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &queue)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return queue, etag, res, fmt.Errorf("failed to fetch skill queue for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return queue, etag, res, nil

}

func characterSkillQueueKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterSkillQueue.String(), strconv.Itoa(int(mods.characterID)))
}

func characterSkillQueuePathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterSkillQueue].Path, mods.characterID)

}
