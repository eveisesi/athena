package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func (s *service) GetUniverseAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error) {

	path := "/v1/universe/ancestries/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &ancestries)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return ancestries, res, fmt.Errorf("failed to fetch ancestries, received status code of %d", sc)
	}

	return ancestries, res, nil

}

func (s *service) GetUniverseBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error) {

	path := "/v1/universe/bloodlines/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &bloodlines)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return bloodlines, res, fmt.Errorf("failed to fetch bloodlines, received status code of %d", sc)
	}

	return bloodlines, res, nil

}

func (s *service) GetUniverseRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error) {

	path := "/v1/universe/races/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &races)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return races, res, fmt.Errorf("failed to fetch races, received status code of %d", sc)
	}

	return races, res, nil

}

func (s *service) GetUniverseFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error) {

	path := "/v2/universe/factions/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &factions)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return factions, res, fmt.Errorf("failed to fetch factions, received status code of %d", sc)
	}

	return factions, res, nil

}

func (s *service) GetUniverseCategories(ctx context.Context, ids []int) ([]int, *http.Response, error) {

	path := "/v1/universe/categories/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &ids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return ids, res, fmt.Errorf("failed to fetch categories, received status code of %d", sc)
	}

	return ids, res, nil

}

func (s *service) GetUniverseCategoriesCategoryID(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error) {

	path := fmt.Sprintf("/v1/universe/categories/%d/", category.CategoryID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &category)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return category, res, fmt.Errorf("failed to fetch category %d, received status code of %d", category.CategoryID, sc)
	}

	return category, res, nil

}

func (s *service) GetUniverseGroupsGroupID(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error) {

	path := fmt.Sprintf("/v1/universe/groups/%d/", group.GroupID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &group)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return group, res, fmt.Errorf("failed to fetch group %d, received status code of %d", group.GroupID, sc)
	}

	return group, res, nil

}

func (s *service) GetUniverseTypesTypeID(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error) {

	path := fmt.Sprintf("/v3/universe/types/%d/", item.TypeID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &item)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return item, res, fmt.Errorf("failed to fetch type %d, received status code of %d", item.TypeID, sc)
	}

	return item, res, nil

}

func (s *service) GetUniverseRegions(ctx context.Context, ids []int) ([]int, *http.Response, error) {

	path := "/v1/universe/regions/"

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &ids)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return ids, res, fmt.Errorf("failed to fetch regions, received status code of %d", sc)
	}

	return ids, res, nil

}

func (s *service) GetUniverseRegionsRegionID(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error) {

	path := fmt.Sprintf("/v1/universe/regions/%d/", region.RegionID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &region)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return region, res, fmt.Errorf("failed to fetch region %d, received status code of %d", region.RegionID, sc)
	}

	return region, res, nil

}

func (s *service) GetUniverseConstellationsConstellationID(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error) {

	path := fmt.Sprintf("/v1/universe/constellations/%d/", constellation.ConstellationID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &constellation)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return constellation, res, fmt.Errorf("failed to fetch constellation %d, received status code of %d", constellation.ConstellationID, sc)
	}

	return constellation, res, nil

}

func (s *service) GetUniverseSolarSystemsSolarSystemID(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error) {

	path := fmt.Sprintf("/v4/universe/systems/%d/", solarSystem.SystemID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &solarSystem)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return solarSystem, res, fmt.Errorf("failed to fetch solarSystem %d, received status code of %d", solarSystem.SystemID, sc)
	}

	return solarSystem, res, nil

}

func (s *service) GetUniverseStationsStationID(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error) {

	path := fmt.Sprintf("/v2/universe/stations/%d/", station.StationID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &station)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return station, res, fmt.Errorf("failed to fetch station %d, received status code of %d", station.SystemID, sc)
	}

	return station, res, nil

}

func (s *service) GetUniverseStructuresStructureID(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error) {

	path := fmt.Sprintf("/v2/universe/structures/%d/", structure.StructureID)

	b, res, err := s.request(ctx, WithMethod(http.MethodGet), WithPath(path), WithAuthorization(member.AccessToken))
	if err != nil {
		return nil, nil, err
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &structure)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return structure, res, fmt.Errorf("failed to fetch structure %d, received status code of %d", structure.StructureID, sc)
	}

	return structure, res, nil

}
