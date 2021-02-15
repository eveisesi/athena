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

func (s *service) GetCharacterAssets(ctx context.Context, member *athena.Member, assets []*athena.MemberAsset) ([]*athena.MemberAsset, *http.Response, error) {

	endpoint := endpoints[GetCharacterAssets]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return assets, res, fmt.Errorf("failed to fetch assets for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {

		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return assets, res, nil

	}

	pages := s.retrieveXPagesFromHeader(res.Header)
	if pages == 0 {
		return nil, nil, fmt.Errorf("received 0 for X-Pages on request %s, expected number greater than 0", path)
	}

	for i := 1; i <= pages; i++ {

		pageAssets := make([]*athena.MemberAsset, 0, 1000)

		mods := s.modifiers(ModWithMember(member), ModWithPage(&i))

		path := endpoint.PathFunc(mods)

		b, res, err := s.request(
			ctx,
			WithMethod(http.MethodGet),
			WithPath(path),
			WithPage(i),
			WithAuthorization(member.AccessToken),
		)
		if err != nil {
			return nil, nil, err
		}

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageAssets)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, err
			}

			assets = append(assets, pageAssets...)

		case sc >= http.StatusBadRequest:
			return assets, res, fmt.Errorf("failed to fetch assets for character %d, received status code of %d", member.ID, sc)
		}

	}

	return assets, res, nil

}

func characterAssetsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterAssets].Path, mods.member.ID)

}

func characterAssetsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(GetCharacterAssets.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}
