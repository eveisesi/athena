package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func isCorporationValid(r *athena.Corporation) bool {
	if r.Name == "" || r.Ticker == "" {
		return false
	}

	return true
}

// GetCorporationsCorporationID makes a HTTP GET Request to the /corporations/{corporation_id} endpoint
// for information about the provided corporation
//
// Documentation: https://esi.evetech.net/ui/#/Corporation/get_corporations_corporation_id
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporationsCorporationID(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *http.Response, error) {

	path := s.endpoints[EndpointGetCorporationsCorporationID](corporation)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path), WithEtag(corporation.Etag))
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

		if etag := s.retrieveEtagHeader(res.Header); etag != "" {
			corporation.Etag = etag
		}

		if !isCorporationValid(corporation) {
			return nil, nil, fmt.Errorf("invalid corporation return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return corporation, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporation.CorporationID, sc)
	}

	corporation.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

	return corporation, res, nil
}

func (s *service) resolveGetCorporationsCorporationIDEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Corporation, received nil")
	}

	var thing *athena.Corporation
	var ok bool

	if thing, ok = obj.(*athena.Corporation); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Corporation, got %T", obj))
	}

	return fmt.Sprintf("/v4/corporations/%d/", thing.CorporationID)

}
