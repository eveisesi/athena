package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type corporationInterface interface {
	GetCorporation(ctx context.Context, corporationID uint) (*athena.Corporation, *athena.Etag, *http.Response, error)
	GetCorporationAllianceHistory(ctx context.Context, corporationID uint) ([]*athena.CorporationAllianceHistory, *athena.Etag, *http.Response, error)
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
func (s *service) GetCorporation(ctx context.Context, corporationID uint) (*athena.Corporation, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCorporation]

	mods := s.modifiers(ModWithCorporationID(corporationID))

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporationID, res.StatusCode)
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

	var corporation = new(athena.Corporation)
	err = json.Unmarshal(b, corporation)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	corporation.ID = corporationID

	if !isCorporationValid(corporation) {
		return nil, nil, nil, fmt.Errorf("invalid corporation return from esi, missing name or ticker")
	}

	return corporation, etag, res, nil
}

func corporationKeyFunc(mods *modifiers) string {

	requireCorporationID(mods)

	return buildKey(GetCorporation.String(), strconv.FormatUint(uint64(mods.corporationID), 10))

}

func corporationPathFunc(mods *modifiers) string {

	requireCorporationID(mods)

	return fmt.Sprintf(endpoints[GetCorporation].Path, mods.corporationID)

}

// GetCorporationAllianceHistory makes a HTTP GET Request to the /v2/corporations/{corporation_id}/alliancehistory/ endpoint
// for information about the provided corporations alliance history
//
// Documentation: https://esi.evetech.net/ui/?version=_latest#/Corporation/get_corporations_corporation_id_alliancehistory
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporationAllianceHistory(ctx context.Context, corporationID uint) ([]*athena.CorporationAllianceHistory, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCorporationAllianceHistory]

	mods := s.modifiers(ModWithCorporationID(corporationID))

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	var history = make([]*athena.CorporationAllianceHistory, 0, 512)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &history)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		etag.Etag = RetrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return history, etag, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporationID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag: %w", err)
	}

	return history, etag, res, nil

}

func corporationAllianceHistoryPathFunc(mods *modifiers) string {

	requireCorporationID(mods)

	return fmt.Sprintf(endpoints[GetCorporationAllianceHistory].Path, mods.corporationID)

}

func corporationAllianceHistoryKeyFunc(mods *modifiers) string {

	requireCorporationID(mods)

	return buildKey(GetCorporationAllianceHistory.String(), strconv.FormatUint(uint64(mods.corporationID), 10))

}
