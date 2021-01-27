package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/athena"
)

func (s *service) GetAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error) {

	path := s.endpoints[EndpointGetAncestries](nil)

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

func (s *service) resolveUniverseAncestriesEndpoint(obj interface{}) string {
	return "/v1/universe/ancestries/"
}

func (s *service) GetBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error) {

	path := s.endpoints[EndpointGetBloodlines](nil)

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

func (s *service) resolveUniverseBloodlinesEndpoint(obj interface{}) string {
	return "/v1/universe/bloodlines/"
}

func (s *service) GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error) {

	path := s.endpoints[EndpointGetRaces](nil)

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

func (s *service) resolveUniverseRacesEndpoint(obj interface{}) string {
	return "/v1/universe/bloodlines/"
}

func (s *service) GetFaction(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error) {

	path := s.endpoints[EndpointGetFaction](nil)

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

func (s *service) resolveUniverseFactionsEndpoint(obj interface{}) string {
	return "/v2/universe/factions/"
}

func (s *service) GetCategories(ctx context.Context, ids []int) ([]int, *http.Response, error) {

	path := s.endpoints[EndpointGetCategories](nil)

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

func (s *service) resolveUniverseCategoriesEndpoint(obj interface{}) string {
	return "/v1/universe/categories/"
}

func (s *service) GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error) {

	path := s.endpoints[EndpointGetCategory](category)

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

func (s *service) resolveUniverseCategoriesCategoryIDEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Category, received nil")
	}

	var thing *athena.Category
	var ok bool
	if thing, ok = obj.(*athena.Category); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Category, got %T", obj))
	}

	return fmt.Sprintf("/v1/universe/categories/%d/", thing.CategoryID)

}

func (s *service) GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error) {

	path := s.endpoints[EndpointGetGroup](group)

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

func (s *service) resolveUniverseGroupsGroupIDEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Group, received nil")
	}

	var thing *athena.Group
	var ok bool

	if thing, ok = obj.(*athena.Group); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Group, got %T", obj))
	}

	return fmt.Sprintf("/v1/universe/groups/%d/", thing.GroupID)

}

func (s *service) GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error) {

	path := s.endpoints[EndpointGetType](item)

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

func (s *service) resolveGetTypeEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Type, received nil")
	}

	var thing *athena.Type
	var ok bool

	if thing, ok = obj.(*athena.Type); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Type, got %T", obj))
	}

	return fmt.Sprintf("/v3/universe/types/%d/", thing.TypeID)

}

func (s *service) GetRegions(ctx context.Context, ids []int) ([]int, *http.Response, error) {

	path := s.endpoints[EndpointGetRegions](nil)

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

func (s *service) resolveGetRegionsEndpoint(obj interface{}) string {

	return "/v1/universe/regions/"

}

func (s *service) GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error) {

	path := s.endpoints[EndpointGetRegion](region)

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

func (s *service) resolveGetRegionEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Region, received nil")
	}

	var thing *athena.Region
	var ok bool

	if thing, ok = obj.(*athena.Region); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Region, got %T", obj))
	}

	return fmt.Sprintf("/v1/universe/regions/%d/", thing.RegionID)

}

func (s *service) GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error) {

	path := s.endpoints[EndpointGetConstellation](constellation)

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

func (s *service) resolveGetConstellationEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Constellation, received nil")
	}

	var thing *athena.Constellation
	var ok bool

	if thing, ok = obj.(*athena.Constellation); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Constellation, got %T", obj))
	}

	return fmt.Sprintf("/v1/universe/constellations/%d/", thing.ConstellationID)

}

func (s *service) GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error) {

	path := s.endpoints[EndpointGetSolarSystem](solarSystem)

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

func (s *service) resolveGetSolarSystemEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.SolarSystem, received nil")
	}

	var thing *athena.SolarSystem
	var ok bool

	if thing, ok = obj.(*athena.SolarSystem); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.SolarSystem, got %T", obj))
	}

	return fmt.Sprintf("/v4/universe/systems/%d/", thing.ConstellationID)

}

func (s *service) GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error) {

	path := s.endpoints[EndpointGetStation](station)

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

func (s *service) resolveGetStationEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Station, received nil")
	}

	var thing *athena.Station
	var ok bool

	if thing, ok = obj.(*athena.Station); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Station, got %T", obj))
	}

	return fmt.Sprintf("/v2/universe/stations/%d/", thing.StationID)

}

func (s *service) GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error) {

	path := s.endpoints[EndpointGetStructure](structure)

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

func (s *service) resolveGetStructureEndpoint(obj interface{}) string {

	if obj == nil {
		panic("invalid type provided for endpoint resolution, expect *athena.Structure, received nil")
	}

	var thing *athena.Structure
	var ok bool

	if thing, ok = obj.(*athena.Structure); !ok {
		panic(fmt.Sprintf("invalid type received for endpoint resolution, expect *athena.Structure, got %T", obj))
	}

	return fmt.Sprintf("/v2/universe/structures/%d/", thing.StructureID)

}
