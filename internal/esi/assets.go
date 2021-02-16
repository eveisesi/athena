package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type assetInterface interface {
	GetCharacterAssets(ctx context.Context, member *athena.Member, assets []*athena.MemberAsset) ([]*athena.MemberAsset, *http.Response, error)
}

func (s *service) HeadCharacterAssets(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterAssets]

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

func (s *service) GetCharacterAssets(ctx context.Context, member *athena.Member, page uint) ([]*athena.MemberAsset, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterAssets]

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	assets := make([]*athena.MemberAsset, 0, 1000) // ESI Specification states a max of 1000 items can be returned per page

	if res.StatusCode >= http.StatusBadRequest {
		return assets, etag, res, fmt.Errorf("failed to fetch assets for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {

		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return assets, etag, res, nil

	}

	err = json.Unmarshal(b, &assets)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return assets, etag, res, nil

}

func characterAssetsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterAssets].Path, mods.member.ID)

}

func characterAssetsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(GetCharacterAssets.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}
