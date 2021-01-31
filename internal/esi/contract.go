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

func (s *service) GetCharacterContractItems(ctx context.Context, member *athena.Member, contracts []*athena.MemberContract) ([]*athena.MemberContract, *http.Response, error) {

	endpoint := s.endpoints[GetCharacterContracts.Name]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return contracts, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.CharacterID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	pages := s.retrieveXPagesFromHeader(res.Header)
	if pages == 0 {
		return nil, nil, fmt.Errorf("received 0 for X-Pages on request %s, expected number greater than 0", path)
	}

	for i := 1; i <= pages; i++ {

		pageContracts := make([]*athena.MemberContract, 0)

		mods := s.modifiers(ModWithMember(member), ModWithPage(&i))

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
			WithPage(i),
			WithAuthorization(member.AccessToken),
		)
		if err != nil {
			return nil, nil, err
		}

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageContracts)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, err
			}

			contracts = append(contracts, pageContracts...)

			etag.Etag = s.retrieveEtagHeader(res.Header)

		case sc >= http.StatusBadRequest:
			return contracts, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.CharacterID, sc)
		}

		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return contracts, res, nil

}

func (s *service) characterContractsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	param := append(make([]string, 0), GetCharacterContracts.Name, strconv.FormatUint(mods.member.CharacterID, 10))

	if mods.page != nil {
		param = append(param, strconv.Itoa(*mods.page))
	}

	return buildKey(param...)
}

func (s *service) characterContractsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContracts.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterContractsEndpoint() *endpoint {
	GetCharacterContracts.KeyFunc = s.characterContractsKeyFunc
	GetCharacterContracts.PathFunc = s.characterContractsPathFunc
	return GetCharacterContracts
}

func (s *service) characterContractItemsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterContractItems.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterContractItemsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContractItems.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterContractItemsEndpoint() *endpoint {
	GetCharacterContractItems.KeyFunc = s.characterContractItemsKeyFunc
	GetCharacterContractItems.PathFunc = s.characterContractItemsPathFunc
	return GetCharacterContractItems
}

func (s *service) characterContractBidsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	return buildKey(GetCharacterContractBids.Name, strconv.FormatUint(mods.member.CharacterID, 10))
}

func (s *service) characterContractBidsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCharacterContractBids.FmtPath, mods.member.CharacterID),
	}

	return u.String()

}

func (s *service) newGetCharacterContractBidsEndpoint() *endpoint {
	GetCharacterContractBids.KeyFunc = s.characterContractBidsKeyFunc
	GetCharacterContractBids.PathFunc = s.characterContractBidsPathFunc
	return GetCharacterContractBids
}
