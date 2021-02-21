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
	HeadCharacterContracts(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error)
	GetCharacterContracts(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberContract, *athena.Etag, *http.Response, error)
	GetCharacterContractItems(ctx context.Context, characterID, contractID uint, token string) ([]*athena.MemberContractItem, *athena.Etag, *http.Response, error)
	GetCharacterContractBids(ctx context.Context, characterID, contractID uint, token string) ([]*athena.MemberContractBid, *athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterContracts(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContracts]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithPage(page))

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
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag: %w", err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterContracts(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberContract, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContracts]

	contracts := make([]*athena.MemberContract, 0)

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithPage(page))

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
		WithPage(page),
		WithAuthorization(token),
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

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return contracts, etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", characterID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	return contracts, etag, res, nil

}

func characterContractsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	params := make([]string, 0, 3)
	params = append(params, GetCharacterContracts.String(), strconv.FormatUint(uint64(mods.characterID), 10))

	if mods.page > 0 {
		params = append(params, strconv.FormatUint(uint64(mods.page), 10))
	}

	return buildKey(params...)

}

func characterContractsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterContracts].Path, mods.characterID)

}

func (s *service) GetCharacterContractItems(ctx context.Context, characterID, contractID uint, token string) ([]*athena.MemberContractItem, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContractItems]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithContractID(contractID))

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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch contract items for character %d contract %d, received status code of %d", characterID, contractID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, etag, res, nil
	}

	var items = make([]*athena.MemberContractItem, 0, 2000)
	err = json.Unmarshal(b, &items)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return items, etag, res, nil

}

func characterContractItemsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireContractID(mods)

	return buildKey(GetCharacterContractItems.String(), strconv.FormatUint(uint64(mods.characterID), 10), strconv.FormatUint(uint64(mods.contractID), 10))
}

func characterContractItemsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireContractID(mods)

	return fmt.Sprintf(endpoints[GetCharacterContractItems].Path, mods.characterID, mods.contractID)

}

func (s *service) GetCharacterContractBids(ctx context.Context, characterID, contractID uint, token string) ([]*athena.MemberContractBid, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterContractBids]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithContractID(contractID))

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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch contract bids for character %d contract %d, received status code of %d", characterID, contractID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, etag, res, nil
	}

	var bids = make([]*athena.MemberContractBid, 0, 2000)
	err = json.Unmarshal(b, &bids)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return bids, etag, res, nil

}

func characterContractBidsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireContractID(mods)

	return buildKey(GetCharacterContractBids.String(), strconv.FormatUint(uint64(mods.characterID), 10), strconv.FormatUint(uint64(mods.contractID), 10))

}

func characterContractBidsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireContractID(mods)

	return fmt.Sprintf(endpoints[GetCharacterContractBids].Path, mods.characterID, mods.contractID)

}
