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

func (s *service) GetAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error) {

	endpoint := s.endpoints[GetAncestries.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) ancestriesKeyFunc(mods *modifiers) string {

	return buildKey(GetAncestries.Name)

}

func (s *service) ancestriesPathFunc(mods *modifiers) string {

	return GetAncestries.FmtPath

}

func (s *service) newGetAncestriesEndpoint() *endpoint {

	GetAncestries.KeyFunc = s.ancestriesKeyFunc
	GetAncestries.PathFunc = s.ancestriesPathFunc
	return GetAncestries

}

func (s *service) GetBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error) {

	endpoint := s.endpoints[GetBloodlines.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) racesKeyFunc(mods *modifiers) string {

	return buildKey(GetRaces.Name)

}

func (s *service) racesPathFunc(mods *modifiers) string {

	return GetRaces.FmtPath

}

func (s *service) newGetRacesEndpoint() *endpoint {

	GetRaces.KeyFunc = s.racesKeyFunc
	GetRaces.PathFunc = s.racesPathFunc
	return GetRaces

}

func (s *service) GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error) {

	endpoint := s.endpoints[GetRaces.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) bloodlinesKeyFunc(mods *modifiers) string {

	return buildKey(GetBloodlines.Name)

}

func (s *service) bloodlinesPathFunc(mods *modifiers) string {

	return GetBloodlines.FmtPath

}

func (s *service) newGetBloodlinesEndpoint() *endpoint {

	GetBloodlines.KeyFunc = s.bloodlinesKeyFunc
	GetBloodlines.PathFunc = s.bloodlinesPathFunc
	return GetBloodlines

}

func (s *service) GetFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error) {

	endpoint := s.endpoints[GetFactions.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) factionsKeyFunc(mods *modifiers) string {

	return buildKey(GetFactions.Name)

}

func (s *service) factionsPathFunc(mods *modifiers) string {

	return GetFactions.FmtPath

}

func (s *service) newGetFactionsEndpoint() *endpoint {

	GetFactions.KeyFunc = s.factionsKeyFunc
	GetFactions.PathFunc = s.factionsPathFunc
	return GetFactions

}

func (s *service) GetCategories(ctx context.Context, ids []uint) ([]uint, *http.Response, error) {

	endpoint := s.endpoints[GetCategories.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) categoriesKeyFunc(mods *modifiers) string {

	return buildKey(GetCategories.Name)

}

func (s *service) categoriesPathFunc(mods *modifiers) string {

	return GetCategories.FmtPath

}

func (s *service) newGetCategoriesEndpoint() *endpoint {

	GetCategories.KeyFunc = s.categoriesKeyFunc
	GetCategories.PathFunc = s.categoriesPathFunc
	return GetCategories

}

func (s *service) GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error) {

	endpoint := s.endpoints[GetCategories.Name]

	mods := s.modifiers(ModWithCategory(category))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &category)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return category, res, fmt.Errorf("failed to fetch category %d, received status code of %d", category.ID, sc)
	}

	return category, res, nil

}

func (s *service) categoryKeyFunc(mods *modifiers) string {

	if mods.category == nil {
		panic("expected type *athena.Category to be provided, received nil for category instead")
	}

	return buildKey(GetCategory.Name, strconv.FormatUint(uint64(mods.category.ID), 10))

}

func (s *service) categoryPathFunc(mods *modifiers) string {

	if mods.category == nil {
		panic("expected type *athena.Category to be provided, received nil for category instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetCategory.FmtPath, mods.category.ID),
	}

	return u.String()

}

func (s *service) newGetCategoryEndpoint() *endpoint {

	GetCategory.KeyFunc = s.categoryKeyFunc
	GetCategory.PathFunc = s.categoryPathFunc
	return GetCategory

}

func (s *service) GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error) {

	endpoint := s.endpoints[GetGroup.Name]

	mods := s.modifiers(ModWithGroup(group))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &group)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return group, res, fmt.Errorf("failed to fetch group %d, received status code of %d", group.ID, sc)
	}

	return group, res, nil

}

func (s *service) groupKeyFunc(mods *modifiers) string {

	if mods.group == nil {
		panic("expected type *athena.Group to be provided, received nil for group instead")
	}

	return buildKey(GetGroup.Name, strconv.FormatUint(uint64(mods.group.ID), 10))

}

func (s *service) groupPathFunc(mods *modifiers) string {

	if mods.group == nil {
		panic("expected type *athena.Group to be provided, received nil for group instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetGroup.FmtPath, mods.group.ID),
	}

	return u.String()

}

func (s *service) newGetGroupEndpoint() *endpoint {

	GetGroup.KeyFunc = s.groupKeyFunc
	GetGroup.PathFunc = s.groupPathFunc
	return GetGroup

}

func (s *service) GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error) {

	endpoint := s.endpoints[GetType.Name]

	mods := s.modifiers(ModWithItem(item))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &item)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return item, res, fmt.Errorf("failed to fetch type %d, received status code of %d", item.ID, sc)
	}

	return item, res, nil

}

func (s *service) typeKeyFunc(mods *modifiers) string {

	if mods.item == nil {
		panic("expected type *athena.Type to be provided, received nil for item instead")
	}

	return buildKey(GetType.Name, strconv.FormatUint(uint64(mods.item.ID), 10))

}

func (s *service) typePathFunc(mods *modifiers) string {

	if mods.item == nil {
		panic("expected type *athena.Type to be provided, received nil for item instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetType.FmtPath, mods.item.ID),
	}

	return u.String()

}

func (s *service) newGetTypeEndpoint() *endpoint {

	GetType.KeyFunc = s.typeKeyFunc
	GetType.PathFunc = s.typePathFunc
	return GetType

}

func (s *service) GetRegions(ctx context.Context, ids []int) ([]int, *http.Response, error) {

	endpoint := s.endpoints[GetRegions.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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

func (s *service) regionsKeyFunc(mods *modifiers) string {

	return buildKey(GetRegions.Name)

}

func (s *service) regionsPathFunc(mods *modifiers) string {

	return GetRegions.FmtPath

}

func (s *service) newGetRegionsEndpoint() *endpoint {

	GetRegions.KeyFunc = s.regionsKeyFunc
	GetRegions.PathFunc = s.regionsPathFunc
	return GetRegions

}

func (s *service) GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error) {

	endpoint := s.endpoints[GetRegion.Name]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &region)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return region, res, fmt.Errorf("failed to fetch region %d, received status code of %d", region.ID, sc)
	}

	return region, res, nil

}

func (s *service) regionKeyFunc(mods *modifiers) string {

	if mods.region == nil {
		panic("expected type *athena.Region to be provided, received nil for region instead")
	}

	return buildKey(GetRegion.Name, strconv.FormatUint(uint64(mods.region.ID), 10))

}

func (s *service) regionPathFunc(mods *modifiers) string {

	if mods.region == nil {
		panic("expected type *athena.Region to be provided, received nil for region instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetRegion.FmtPath, mods.region.ID),
	}

	return u.String()

}

func (s *service) newGetRegionEndpoint() *endpoint {

	GetRegion.KeyFunc = s.regionKeyFunc
	GetRegion.PathFunc = s.regionPathFunc
	return GetRegion

}

func (s *service) GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error) {

	endpoint := s.endpoints[GetRaces.Name]

	mods := s.modifiers(ModWithConstellation(constellation))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &constellation)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return constellation, res, fmt.Errorf("failed to fetch constellation %d, received status code of %d", constellation.ID, sc)
	}

	return constellation, res, nil

}

func (s *service) constellationKeyFunc(mods *modifiers) string {

	if mods.constellation == nil {
		panic("expected type *athena.Constellation to be provided, received nil for constellation instead")
	}

	return buildKey(GetConstellation.Name, strconv.FormatUint(uint64(mods.constellation.ID), 10))

}

func (s *service) constellationPathFunc(mods *modifiers) string {

	if mods.constellation == nil {
		panic("expected type *athena.Constellation to be provided, received nil for constellation instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetConstellation.FmtPath, mods.constellation.ID),
	}

	return u.String()

}

func (s *service) newGetConstellationEndpoint() *endpoint {

	GetConstellation.KeyFunc = s.constellationKeyFunc
	GetConstellation.PathFunc = s.constellationPathFunc
	return GetConstellation

}

func (s *service) GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error) {

	endpoint := s.endpoints[GetRaces.Name]

	mods := s.modifiers(ModWithSystem(solarSystem))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &solarSystem)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return solarSystem, res, fmt.Errorf("failed to fetch solarSystem %d, received status code of %d", solarSystem.ID, sc)
	}

	return solarSystem, res, nil

}

func (s *service) solarSystemKeyFunc(mods *modifiers) string {

	if mods.solarSystem == nil {
		panic("expected type *athena.SolarSystem to be provided, received nil for solarSystem instead")
	}

	return buildKey(GetSolarSystem.Name, strconv.FormatUint(uint64(mods.solarSystem.ID), 10))

}

func (s *service) solarSystemPathFunc(mods *modifiers) string {

	if mods.solarSystem == nil {
		panic("expected type *athena.SolarSystem to be provided, received nil for solarSystem instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetSolarSystem.FmtPath, mods.solarSystem.ID),
	}

	return u.String()

}

func (s *service) newGetSolarSystemEndpoint() *endpoint {

	GetSolarSystem.KeyFunc = s.solarSystemKeyFunc
	GetSolarSystem.PathFunc = s.solarSystemPathFunc
	return GetSolarSystem

}

func (s *service) GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error) {

	endpoint := s.endpoints[GetStation.Name]

	mods := s.modifiers(ModWithStation(station))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		err = json.Unmarshal(b, &station)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return station, res, fmt.Errorf("failed to fetch station %d, received status code of %d", station.ID, sc)
	}

	return station, res, nil

}

func (s *service) stationKeyFunc(mods *modifiers) string {

	if mods.station == nil {
		panic("expected type *athena.Station to be provided, received nil for station instead")
	}

	return buildKey(GetStation.Name, strconv.FormatUint(uint64(mods.station.ID), 10))

}

func (s *service) stationPathFunc(mods *modifiers) string {

	if mods.station == nil {
		panic("expected type *athena.Station to be provided, received nil for station instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetStation.FmtPath, mods.station.ID),
	}

	return u.String()

}

func (s *service) newGetStationEndpoint() *endpoint {

	GetStation.KeyFunc = s.stationKeyFunc
	GetStation.PathFunc = s.stationPathFunc
	return GetStation

}

func (s *service) GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error) {

	endpoint := s.endpoints[GetStructure.Name]

	mods := s.modifiers(ModWithStructure(structure))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
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
		return structure, res, fmt.Errorf("failed to fetch structure %d, received status code of %d", structure.ID, sc)
	}

	return structure, res, nil

}

func (s *service) structureKeyFunc(mods *modifiers) string {

	if mods.structure == nil {
		panic("expected type *athena.Structure to be provided, received nil for structure instead")
	}

	return buildKey(GetStation.Name, strconv.FormatUint(mods.structure.ID, 10))

}

func (s *service) structurePathFunc(mods *modifiers) string {

	if mods.structure == nil {
		panic("expected type *athena.Structure to be provided, received nil for structure instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(GetStructure.FmtPath, mods.structure.ID),
	}

	return u.String()

}

func (s *service) newGetStructureEndpoint() *endpoint {

	GetStructure.KeyFunc = s.structureKeyFunc
	GetStructure.PathFunc = s.structurePathFunc
	return GetStructure

}
