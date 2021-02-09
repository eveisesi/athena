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

func isCorporationValid(r *athena.Corporation) bool {
	if r.Name == "" || r.Ticker == "" {
		return false
	}

	return true
}

// GetCorporation makes a HTTP GET Request to the /corporations/{corporation_id} endpoint
// for information about the provided corporation
//
// Documentation: https://esi.evetech.net/ui/#/Corporation/get_corporations_corporation_id
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *http.Response, error) {

	endpoint := s.endpoints[GetCorporation.Name]

	mods := s.modifiers(ModWithCorporation(corporation))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	// WithEtag(etag)
	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
	)
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, corporation)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

		if !isCorporationValid(corporation) {
			return nil, nil, fmt.Errorf("invalid corporation return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return corporation, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporation.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return corporation, res, nil
}

func (s *service) corporationKeyFunc(mods *modifiers) string {

	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	return buildKey(GetCorporation.Name, strconv.Itoa(int(mods.corporation.ID)))
}

func (s *service) corporationPathFunc(mods *modifiers) string {

	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCorporation.FmtPath, mods.corporation.ID),
	}

	return u.String()

}

func (s *service) newGetCorporationEndpoint() *endpoint {
	GetCorporation.KeyFunc = s.corporationKeyFunc
	GetCorporation.PathFunc = s.corporationPathFunc
	return GetCorporation
}

// GetCorporationAllianceHistory makes a HTTP GET Request to the /v2/corporations/{corporation_id}/alliancehistory/ endpoint
// for information about the provided corporations alliance history
//
// Documentation: https://esi.evetech.net/ui/?version=_latest#/Corporation/get_corporations_corporation_id_alliancehistory
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporationAllianceHistory(ctx context.Context, corporation *athena.Corporation, history []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, *athena.Etag, *http.Response, error) {

	endpoint := s.endpoints[GetCorporationAllianceHistory.Name]

	mods := s.modifiers(ModWithCorporation(corporation))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	// WithEtag(etag)
	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &history)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return history, etag, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporation.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return history, etag, res, nil

}

func (s *service) corporationAllianceHistoryPathFunc(mods *modifiers) string {
	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCorporationAllianceHistory.FmtPath, mods.corporation.ID),
	}

	return u.String()
}

func (s *service) corporationAllianceHistoryKeyFunc(mods *modifiers) string {
	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	return buildKey(GetCorporationAllianceHistory.Name, strconv.Itoa(int(mods.corporation.ID)))
}

func (s *service) newGetCorporationAllianceHistoryEndpoint() *endpoint {
	GetCorporationAllianceHistory.KeyFunc = s.corporationAllianceHistoryKeyFunc
	GetCorporationAllianceHistory.PathFunc = s.corporationAllianceHistoryPathFunc
	return GetCorporationAllianceHistory
}
