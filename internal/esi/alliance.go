package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

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
func (s *service) GetAlliance(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, *http.Response, error) {

	path := s.endpoints[EndpointGetAlliance](alliance)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path), WithEtag(alliance.Etag))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, alliance)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			alliance.Etag = etag
		}

		if !isAllianceValid(alliance) {
			return nil, nil, fmt.Errorf("invalid alliance return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return alliance, res, fmt.Errorf("failed to fetch alliance %d, received status code of %d", alliance.AllianceID, sc)
	}

	alliance.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return alliance, res, nil
}

func (s *service) resolveGetAllianceEndpoint(obj interface{}) string {

	mods := s.modifiers(modFuncs)

	if mods

	var thing *athena.Alliance
	var ok bool

	if thing, ok = obj.(*athena.Alliance); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Alliance, got %T", obj))
	}

	return fmt.Sprintf("/v3/alliances/%d/", thing.AllianceID)

}
