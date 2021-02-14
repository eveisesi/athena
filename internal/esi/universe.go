package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type universeInterface interface {
	GetAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error)
	GetBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error)
	GetCategories(ctx context.Context, ids []uint) ([]uint, *http.Response, error)
	GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error)
	GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error)
	GetFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error)
	GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error)
	GetMoon(ctx context.Context, moon *athena.Moon) (*athena.Moon, *http.Response, error)
	GetAsteroidBelt(ctx context.Context, belt *athena.AsteroidBelt) (*athena.AsteroidBelt, *http.Response, error)
	GetPlanet(ctx context.Context, planet *athena.Planet) (*athena.Planet, *http.Response, error)
	GetRegions(ctx context.Context, ids []uint) ([]uint, *http.Response, error)
	GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error)
	GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error)
	GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error)
	GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error)
	GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error)
	GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error)

	PostUniverseNames(ctx context.Context, ids []uint) ([]*PostUniverseNamesOK, *http.Response, error)
}

type PostUniverseNamesOK struct {
	Category Category `json:"category"`
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
}

type Category string

const (
	CategoryAlliance      Category = "alliance"
	CategoryCharacter     Category = "character"
	CategoryCorporation   Category = "corporation"
	CategoryConstellation Category = "constellation"
	CategoryInventoryType Category = "inventory_type"
	CategoryRegion        Category = "region"
	CategorySolarSystem   Category = "solar_system"
	CategoryStation       Category = "station"
	CategoryFaction       Category = "faction"
)

func (s *service) PostUniverseNames(ctx context.Context, ids []uint) ([]*PostUniverseNamesOK, *http.Response, error) {

	endpoint := endpoints[PostUniverseNames]

	mods := s.modifiers()

	path := endpoint.PathFunc(mods)

	if len(ids) == 0 {
		return nil, nil, fmt.Errorf("0 ids received, must be greater than 0")
	}

	if len(ids) > 250 {
		return nil, nil, fmt.Errorf("more than 250 ids received, limit is 250")
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marsahl ids to byte array: %w", err)
	}

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodPost),
		WithPath(path),
		WithBody(data),
	)
	if err != nil {
		return nil, nil, err
	}

	var names = make([]*PostUniverseNamesOK, 0)

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &names)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		}

	case sc >= http.StatusBadRequest:
		var e = new(GenericError)
		err = json.Unmarshal(b, e)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to unmarsahl error onto error struct, data: %s: %w", string(b), err)
		}

		return nil, res, fmt.Errorf("failed to post names, received status code of %d: %w", sc, e)
	}

	return names, res, nil

}

func universeNamesPathFunc(mods *modifiers) string {
	return endpoints[PostUniverseNames].Path
}

func (s *service) GetAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error) {

	endpoint := endpoints[GetAncestries]

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

func ancestriesKeyFunc(mods *modifiers) string {

	return buildKey(GetAncestries.String())

}

func ancestriesPathFunc(mods *modifiers) string {

	return endpoints[GetAncestries].Path

}

func (s *service) GetBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error) {

	endpoint := endpoints[GetBloodlines]

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

func bloodlinesKeyFunc(mods *modifiers) string {

	return buildKey(GetBloodlines.String())

}

func bloodlinesPathFunc(mods *modifiers) string {

	return endpoints[GetBloodlines].Path

}

func (s *service) GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error) {

	endpoint := endpoints[GetRaces]

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

func racesKeyFunc(mods *modifiers) string {

	return buildKey(GetRaces.String())

}

func racesPathFunc(mods *modifiers) string {

	return endpoints[GetRaces].Path

}

func (s *service) GetFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error) {

	endpoint := endpoints[GetFactions]

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

func factionsKeyFunc(mods *modifiers) string {

	return buildKey(GetFactions.String())

}

func factionsPathFunc(mods *modifiers) string {

	return endpoints[GetFactions].Path

}

func (s *service) GetCategories(ctx context.Context, ids []uint) ([]uint, *http.Response, error) {

	endpoint := endpoints[GetCategories]

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

func categoriesKeyFunc(mods *modifiers) string {

	return buildKey(GetCategories.String())

}

func categoriesPathFunc(mods *modifiers) string {

	return endpoints[GetCategories].Path

}

func (s *service) GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error) {

	endpoint := endpoints[GetCategory]

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

func categoryKeyFunc(mods *modifiers) string {

	requireCategory(mods)

	return buildKey(GetCategory.String(), strconv.FormatUint(uint64(mods.category.ID), 10))

}

func categoryPathFunc(mods *modifiers) string {

	requireCategory(mods)

	return fmt.Sprintf(endpoints[GetCategory].Path, mods.category.ID)

}

func (s *service) GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error) {

	endpoint := endpoints[GetGroup]

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

func groupKeyFunc(mods *modifiers) string {

	requireGroup(mods)

	return buildKey(GetGroup.String(), strconv.FormatUint(uint64(mods.group.ID), 10))

}

func groupPathFunc(mods *modifiers) string {

	requireGroup(mods)

	return fmt.Sprintf(endpoints[GetGroup].Path, mods.member.ID)

}

func (s *service) GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error) {

	endpoint := endpoints[GetType]

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

func typeKeyFunc(mods *modifiers) string {

	requireItem(mods)

	return buildKey(GetType.String(), strconv.FormatUint(uint64(mods.item.ID), 10))

}

func typePathFunc(mods *modifiers) string {

	requireItem(mods)

	return fmt.Sprintf(endpoints[GetType].Path, mods.item.ID)

}

func (s *service) GetRegions(ctx context.Context, ids []uint) ([]uint, *http.Response, error) {

	endpoint := endpoints[GetRegions]

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

func regionsKeyFunc(mods *modifiers) string {

	return buildKey(GetRegions.String())

}

func regionsPathFunc(mods *modifiers) string {

	return endpoints[GetRegions].Path

}

func (s *service) GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error) {

	endpoint := endpoints[GetRegion]

	mods := s.modifiers(ModWithRegion(region))

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

func regionKeyFunc(mods *modifiers) string {

	requireRegion(mods)

	return buildKey(GetRegion.String(), strconv.FormatUint(uint64(mods.region.ID), 10))

}

func regionPathFunc(mods *modifiers) string {

	requireRegion(mods)

	return fmt.Sprintf(endpoints[GetRegion].Path, mods.region.ID)

}

func (s *service) GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error) {

	endpoint := endpoints[GetConstellation]

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

func constellationKeyFunc(mods *modifiers) string {

	requireConstellation(mods)

	return buildKey(GetConstellation.String(), strconv.FormatUint(uint64(mods.constellation.ID), 10))

}

func constellationPathFunc(mods *modifiers) string {

	requireConstellation(mods)

	return fmt.Sprintf(endpoints[GetConstellation].Path, mods.constellation.ID)

}

func (s *service) GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error) {

	endpoint := endpoints[GetSolarSystem]

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

func solarSystemKeyFunc(mods *modifiers) string {

	requireSystem(mods)

	return buildKey(GetSolarSystem.String(), strconv.FormatUint(uint64(mods.solarSystem.ID), 10))

}

func solarSystemPathFunc(mods *modifiers) string {

	requireSystem(mods)

	return fmt.Sprintf(endpoints[GetSolarSystem].Path, mods.solarSystem.ID)

}

func (s *service) GetPlanet(ctx context.Context, planet *athena.Planet) (*athena.Planet, *http.Response, error) {

	endpoint := endpoints[GetPlanet]

	mods := s.modifiers(ModWithPlanet(planet))

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
		err = json.Unmarshal(b, &planet)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return planet, res, fmt.Errorf("failed to fetch planet %d, received status code of %d", planet.ID, sc)
	}

	return planet, res, nil

}

func planetKeyFunc(mods *modifiers) string {

	requirePlanet(mods)

	return buildKey(GetPlanet.String(), strconv.FormatUint(uint64(mods.planet.ID), 10))

}

func planetPathFunc(mods *modifiers) string {

	requirePlanet(mods)

	return fmt.Sprintf(endpoints[GetPlanet].Path, mods.planet.ID)

}

func (s *service) GetMoon(ctx context.Context, moon *athena.Moon) (*athena.Moon, *http.Response, error) {

	endpoint := endpoints[GetMoon]

	mods := s.modifiers(ModWithMoon(moon))

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
		err = json.Unmarshal(b, &moon)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return moon, res, fmt.Errorf("failed to fetch moon %d, received status code of %d", moon.ID, sc)
	}

	return moon, res, nil

}

func moonKeyFunc(mods *modifiers) string {

	requireMoon(mods)

	return buildKey(GetMoon.String(), strconv.FormatUint(uint64(mods.moon.ID), 10))

}

func moonPathFunc(mods *modifiers) string {

	requireMoon(mods)

	return fmt.Sprintf(endpoints[GetMoon].Path, mods.moon.ID)

}

func (s *service) GetAsteroidBelt(ctx context.Context, asteroidBelt *athena.AsteroidBelt) (*athena.AsteroidBelt, *http.Response, error) {

	endpoint := endpoints[GetAsteroidBelt]

	mods := s.modifiers(ModWithAsteroidBelt(asteroidBelt))

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
		err = json.Unmarshal(b, &asteroidBelt)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, nil, err
		}

	case sc >= http.StatusBadRequest:
		return asteroidBelt, res, fmt.Errorf("failed to fetch asteroidBelt %d, received status code of %d", asteroidBelt.ID, sc)
	}

	return asteroidBelt, res, nil

}

func beltKeyFunc(mods *modifiers) string {

	requireAsteriodBelt(mods)

	return buildKey(GetAsteroidBelt.String(), strconv.FormatUint(uint64(mods.asteroidBelt.ID), 10))

}

func beltPathFunc(mods *modifiers) string {

	requireAsteriodBelt(mods)

	return fmt.Sprintf(endpoints[GetAsteroidBelt].Path, mods.asteroidBelt.ID)

}

func (s *service) GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error) {

	endpoint := endpoints[GetStation]

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

func stationKeyFunc(mods *modifiers) string {

	requireStation(mods)

	return buildKey(GetStation.String(), strconv.FormatUint(uint64(mods.station.ID), 10))

}

func stationPathFunc(mods *modifiers) string {

	requireStation(mods)

	return fmt.Sprintf(endpoints[GetStation].Path, mods.station.ID)

}

func (s *service) GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error) {

	endpoint := endpoints[GetStructure]

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

func structureKeyFunc(mods *modifiers) string {

	requireStructure(mods)

	return buildKey(GetStructure.String(), strconv.FormatUint(mods.structure.ID, 10))

}

func structurePathFunc(mods *modifiers) string {

	requireStructure(mods)

	return fmt.Sprintf(endpoints[GetStructure].Path, mods.structure.ID)

}
