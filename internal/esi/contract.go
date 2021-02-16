package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type contractInterface interface {
	HeadCharacterContracts(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterContracts(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContracts]

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

func (s *service) GetCharacterContracts(ctx context.Context, member *athena.Member, page uint) ([]*athena.MemberContract, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContracts]

	contracts := make([]*athena.MemberContract, 0)

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

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
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &contracts)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return contracts, etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	return contracts, etag, res, nil

}

func characterContractsKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requirePage(mods)

	return buildKey(GetCharacterContracts.String(), strconv.FormatUint(uint64(mods.member.ID), 10), strconv.FormatUint(uint64(mods.page), 10))

}

func characterContractsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterContracts].Path, mods.member.ID)

}

func (s *service) GetCharacterContractItems(ctx context.Context, member *athena.Member, contract *athena.MemberContract, items []*athena.MemberContractItem) ([]*athena.MemberContractItem, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContractItems]

	mods := s.modifiers(ModWithMember(member), ModWithContract(contract))

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
		err = json.Unmarshal(b, &items)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return items, etag, res, fmt.Errorf("failed to fetch contract items for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
	}

	return items, etag, res, nil

}

func characterContractItemsKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requireContract(mods)

	return buildKey(GetCharacterContractItems.String(), strconv.FormatUint(uint64(mods.member.ID), 10), strconv.FormatUint(uint64(mods.contract.ContractID), 10))
}

func characterContractItemsPathFunc(mods *modifiers) string {

	requireMember(mods)
	requireContract(mods)

	return fmt.Sprintf(endpoints[GetCharacterContractItems].Path, mods.member.ID, mods.contract.ContractID)

}

func (s *service) GetCharacterContractBids(ctx context.Context, member *athena.Member, contract *athena.MemberContract, bids []*athena.MemberContractBid) ([]*athena.MemberContractBid, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContractBids]

	mods := s.modifiers(ModWithMember(member), ModWithContract(contract))

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
		err = json.Unmarshal(b, &bids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return bids, etag, res, fmt.Errorf("failed to fetch contract bids for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
	}

	return bids, etag, res, nil

}

func characterContractBidsKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requireContract(mods)

	return buildKey(GetCharacterContractBids.String(), strconv.FormatUint(uint64(mods.member.ID), 10), strconv.FormatUint(uint64(mods.contract.ContractID), 10))

}

func characterContractBidsPathFunc(mods *modifiers) string {

	requireMember(mods)
	requireContract(mods)

	return fmt.Sprintf(endpoints[GetCharacterContractBids].Path, mods.member.ID, mods.contract.ContractID)

}
