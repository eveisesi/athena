package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type allianceInterface interface {
	GetAlliance(ctx context.Context, allianceID uint) (*athena.Alliance, *athena.Etag, *http.Response, error)
}

func isAllianceValid(r *athena.Alliance) bool {
	if r.Name == "" || r.Ticker == "" {
		return false
	}
	return true
}

// GetAlliance makes a HTTP GET Request to the /alliances/{alliance_id} endpoint
// for information about the provided alliance
//
// Documentation: https://esi.evetech.net/ui/#/Alliance/get_alliances_alliance_id
// Version: v3
// Cache: 3600 sec (1 Hour)
func (s *service) GetAlliance(ctx context.Context, allianceID uint) (*athena.Alliance, *athena.Etag, *http.Response, error) {

	// Fetch configuration for this endpoint
	endpoint := endpoints[GetAlliance]

	// Prime modifiers with alliance
	mods := s.modifiers(ModWithAllianceID(allianceID))

	// Fetch Etag for request
	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	var alliance = new(athena.Alliance)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, alliance)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, nil, err
		}

		alliance.ID = allianceID

		etag.Etag = RetrieveEtagHeader(res.Header)

		if !isAllianceValid(alliance) {
			return nil, nil, nil, fmt.Errorf("invalid alliance returned from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return alliance, etag, res, fmt.Errorf("failed to fetch alliance %d, received status code of %d", alliance.ID, sc)
	}

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return alliance, etag, res, nil

}

func allianceKeyFunc(mods *modifiers) string {

	requireAllianceID(mods)

	return buildKey(GetAlliance.String(), strconv.Itoa(int(mods.allianceID)))

}

func alliancePathFunc(mods *modifiers) string {

	requireAllianceID(mods)

	return fmt.Sprintf(endpoints[GetAlliance].Path, mods.allianceID)

}
