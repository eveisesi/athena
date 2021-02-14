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

type corporationInterface interface {
	GetCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *athena.Etag, *http.Response, error)
	GetCorporationAllianceHistory(ctx context.Context, corporation *athena.Corporation, history []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, *athena.Etag, *http.Response, error)
}

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
func (s *service) GetCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCorporation]

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
		err = json.Unmarshal(b, corporation)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

		if !isCorporationValid(corporation) {
			return nil, nil, nil, fmt.Errorf("invalid corporation return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return corporation, etag, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporation.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return corporation, etag, res, nil
}

func corporationKeyFunc(mods *modifiers) string {

	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	return buildKey(GetCorporation.String(), strconv.Itoa(int(mods.corporation.ID)))
}

func corporationPathFunc(mods *modifiers) string {

	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCorporation].Path, mods.corporation.ID),
	}

	return u.String()

}

// GetCorporationAllianceHistory makes a HTTP GET Request to the /v2/corporations/{corporation_id}/alliancehistory/ endpoint
// for information about the provided corporations alliance history
//
// Documentation: https://esi.evetech.net/ui/?version=_latest#/Corporation/get_corporations_corporation_id_alliancehistory
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporationAllianceHistory(ctx context.Context, corporation *athena.Corporation, history []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCorporationAllianceHistory]

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

func corporationAllianceHistoryPathFunc(mods *modifiers) string {
	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCorporationAllianceHistory].Path, mods.corporation.ID),
	}

	return u.String()
}

func corporationAllianceHistoryKeyFunc(mods *modifiers) string {
	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil for corporation instead")
	}

	return buildKey(GetCorporationAllianceHistory.String(), strconv.Itoa(int(mods.corporation.ID)))
}
