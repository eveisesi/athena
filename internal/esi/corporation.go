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

	// etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	// if err != nil {
	// 	return nil, nil, err
	// }

	path := endpoint.PathFunc(mods)

	// WithEtag(etag)
	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
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

		// etag.Etag = s.retrieveEtagHeader(res.Header)

		if !isCorporationValid(corporation) {
			return nil, nil, fmt.Errorf("invalid corporation return from esi, missing name or ticker")
		}
	case sc >= http.StatusBadRequest:
		return corporation, res, fmt.Errorf("failed to fetch corporation %d, received status code of %d", corporation.ID, sc)
	}

	// etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)

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
