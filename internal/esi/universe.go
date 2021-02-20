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
	GetAncestries(ctx context.Context) ([]*athena.Ancestry, *http.Response, error)
	GetBloodlines(ctx context.Context) ([]*athena.Bloodline, *http.Response, error)
	GetCategories(ctx context.Context) ([]uint, *http.Response, error)
	GetCategory(ctx context.Context, categoryID uint) (*athena.Category, *http.Response, error)
	GetConstellation(ctx context.Context, constellationID uint) (*athena.Constellation, *http.Response, error)
	GetFactions(ctx context.Context) ([]*athena.Faction, *http.Response, error)
	GetGroup(ctx context.Context, groupID uint) (*athena.Group, *http.Response, error)
	GetRegions(ctx context.Context) ([]uint, *http.Response, error)
	GetRegion(ctx context.Context, regionID uint) (*athena.Region, *http.Response, error)
	GetRaces(ctx context.Context) ([]*athena.Race, *http.Response, error)
	GetSolarSystem(ctx context.Context, solarSystemID uint) (*athena.SolarSystem, *http.Response, error)
	GetStation(ctx context.Context, stationID uint) (*athena.Station, *http.Response, error)
	GetStructure(ctx context.Context, structureID uint64, token string) (*athena.Structure, *http.Response, error)
	GetType(ctx context.Context, typeID uint) (*athena.Type, *http.Response, error)

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

func (s *service) GetAncestries(ctx context.Context) ([]*athena.Ancestry, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var ancestries = make([]*athena.Ancestry, 0)

	if res.StatusCode > http.StatusOK {
		return ancestries, res, fmt.Errorf("failed to fetch ancestries, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &ancestries)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return ancestries, res, nil

}

func ancestriesKeyFunc(mods *modifiers) string {
	return buildKey(GetAncestries.String())
}

func ancestriesPathFunc(mods *modifiers) string {
	return endpoints[GetAncestries].Path
}

func (s *service) GetBloodlines(ctx context.Context) ([]*athena.Bloodline, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var bloodlines = make([]*athena.Bloodline, 0)
	if res.StatusCode > http.StatusOK {
		return bloodlines, res, fmt.Errorf("failed to fetch bloodlines, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &bloodlines)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return bloodlines, res, nil

}

func bloodlinesKeyFunc(mods *modifiers) string {
	return buildKey(GetBloodlines.String())
}

func bloodlinesPathFunc(mods *modifiers) string {
	return endpoints[GetBloodlines].Path
}

func (s *service) GetRaces(ctx context.Context) ([]*athena.Race, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var races = make([]*athena.Race, 0)
	if res.StatusCode > http.StatusOK {
		return races, res, fmt.Errorf("failed to fetch races, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &races)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return races, res, nil

}

func racesKeyFunc(mods *modifiers) string {
	return buildKey(GetRaces.String())
}

func racesPathFunc(mods *modifiers) string {
	return endpoints[GetRaces].Path
}

func (s *service) GetFactions(ctx context.Context) ([]*athena.Faction, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var factions = make([]*athena.Faction, 0)
	if res.StatusCode > http.StatusOK {
		return factions, res, fmt.Errorf("failed to fetch factions, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &factions)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return factions, res, nil

}

func factionsKeyFunc(mods *modifiers) string {
	return buildKey(GetFactions.String())
}

func factionsPathFunc(mods *modifiers) string {
	return endpoints[GetFactions].Path
}

func (s *service) GetCategories(ctx context.Context) ([]uint, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var categories = make([]uint, 0)
	if res.StatusCode > http.StatusOK {
		return categories, res, fmt.Errorf("failed to fetch categories, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &categories)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return categories, res, nil

}

func categoriesKeyFunc(mods *modifiers) string {
	return buildKey(GetCategories.String())
}

func categoriesPathFunc(mods *modifiers) string {
	return endpoints[GetCategories].Path
}

func (s *service) GetCategory(ctx context.Context, categoryID uint) (*athena.Category, *http.Response, error) {

	endpoint := endpoints[GetCategory]

	mods := s.modifiers(ModWithCategoryID(categoryID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var category = new(athena.Category)
	if res.StatusCode > http.StatusOK {
		return category, res, fmt.Errorf("failed to fetch category, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, category)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return category, res, nil

}

func categoryKeyFunc(mods *modifiers) string {

	requireCategoryID(mods)

	return buildKey(GetCategory.String(), strconv.FormatUint(uint64(mods.categoryID), 10))

}

func categoryPathFunc(mods *modifiers) string {

	requireCategoryID(mods)

	return fmt.Sprintf(endpoints[GetCategory].Path, mods.categoryID)

}

func (s *service) GetGroup(ctx context.Context, groupID uint) (*athena.Group, *http.Response, error) {

	endpoint := endpoints[GetGroup]

	mods := s.modifiers(ModWithGroupID(groupID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var group = new(athena.Group)
	if res.StatusCode > http.StatusOK {
		return group, res, fmt.Errorf("failed to fetch group, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, group)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return group, res, nil

}

func groupKeyFunc(mods *modifiers) string {

	requireGroupID(mods)

	return buildKey(GetGroup.String(), strconv.FormatUint(uint64(mods.groupID), 10))

}

func groupPathFunc(mods *modifiers) string {

	requireGroupID(mods)

	return fmt.Sprintf(endpoints[GetGroup].Path, mods.groupID)

}

func (s *service) GetType(ctx context.Context, typeID uint) (*athena.Type, *http.Response, error) {

	endpoint := endpoints[GetType]

	mods := s.modifiers(ModWithItemID(typeID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var item = new(athena.Type)
	if res.StatusCode > http.StatusOK {
		return item, res, fmt.Errorf("failed to fetch item, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, item)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return item, res, nil

}

func typeKeyFunc(mods *modifiers) string {

	requireItemID(mods)

	return buildKey(GetType.String(), strconv.FormatUint(uint64(mods.itemID), 10))

}

func typePathFunc(mods *modifiers) string {

	requireItemID(mods)

	return fmt.Sprintf(endpoints[GetType].Path, mods.itemID)

}

func (s *service) GetRegions(ctx context.Context) ([]uint, *http.Response, error) {

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
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var ids = make([]uint, 0)
	if res.StatusCode > http.StatusOK {
		return ids, res, fmt.Errorf("failed to fetch item, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, &ids)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return ids, res, nil

}

func regionsKeyFunc(mods *modifiers) string {
	return buildKey(GetRegions.String())
}

func regionsPathFunc(mods *modifiers) string {
	return endpoints[GetRegions].Path
}

func (s *service) GetRegion(ctx context.Context, regionID uint) (*athena.Region, *http.Response, error) {

	endpoint := endpoints[GetRegion]

	mods := s.modifiers(ModWithRegionID(regionID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var region = new(athena.Region)
	if res.StatusCode > http.StatusOK {
		return region, res, fmt.Errorf("failed to fetch region, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, region)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return region, res, nil

}

func regionKeyFunc(mods *modifiers) string {

	requireRegionID(mods)

	return buildKey(GetRegion.String(), strconv.FormatUint(uint64(mods.regionID), 10))

}

func regionPathFunc(mods *modifiers) string {

	requireRegionID(mods)

	return fmt.Sprintf(endpoints[GetRegion].Path, mods.regionID)

}

func (s *service) GetConstellation(ctx context.Context, constellationID uint) (*athena.Constellation, *http.Response, error) {

	endpoint := endpoints[GetConstellation]

	mods := s.modifiers(ModWithConstellationID(constellationID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var constellation = new(athena.Constellation)
	if res.StatusCode > http.StatusOK {
		return constellation, res, fmt.Errorf("failed to fetch constellation, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, constellation)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return constellation, res, nil

}

func constellationKeyFunc(mods *modifiers) string {

	requireConstellationID(mods)

	return buildKey(GetConstellation.String(), strconv.FormatUint(uint64(mods.constellationID), 10))

}

func constellationPathFunc(mods *modifiers) string {

	requireConstellationID(mods)

	return fmt.Sprintf(endpoints[GetConstellation].Path, mods.constellationID)

}

func (s *service) GetSolarSystem(ctx context.Context, systemID uint) (*athena.SolarSystem, *http.Response, error) {

	endpoint := endpoints[GetSolarSystem]

	mods := s.modifiers(ModWithSystemID(systemID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var system = new(athena.SolarSystem)
	if res.StatusCode > http.StatusOK {
		return system, res, fmt.Errorf("failed to fetch system, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, system)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return system, res, nil

}

func solarSystemKeyFunc(mods *modifiers) string {

	requireSystemID(mods)

	return buildKey(GetSolarSystem.String(), strconv.FormatUint(uint64(mods.solarSystemID), 10))

}

func solarSystemPathFunc(mods *modifiers) string {

	requireSystemID(mods)

	return fmt.Sprintf(endpoints[GetSolarSystem].Path, mods.solarSystemID)

}

func (s *service) GetStation(ctx context.Context, stationID uint) (*athena.Station, *http.Response, error) {

	endpoint := endpoints[GetStation]

	mods := s.modifiers(ModWithStationID(stationID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
	)
	if err != nil {
		return nil, nil, err
	}

	var station = new(athena.Station)
	if res.StatusCode > http.StatusOK {
		return station, res, fmt.Errorf("failed to fetch station, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, station)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return station, res, nil

}

func stationKeyFunc(mods *modifiers) string {

	requireStationID(mods)

	return buildKey(GetStation.String(), strconv.FormatUint(uint64(mods.stationID), 10))

}

func stationPathFunc(mods *modifiers) string {

	requireStationID(mods)

	return fmt.Sprintf(endpoints[GetStation].Path, mods.stationID)

}

func (s *service) GetStructure(ctx context.Context, structureID uint64, token string) (*athena.Structure, *http.Response, error) {

	endpoint := endpoints[GetStructure]

	mods := s.modifiers(ModWithStructureID(structureID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	var structure = new(athena.Structure)
	if res.StatusCode > http.StatusOK {
		return structure, res, fmt.Errorf("failed to fetch structure, received status code of %d", res.StatusCode)
	}

	err = json.Unmarshal(b, structure)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, err
	}

	return structure, res, nil

}

func structureKeyFunc(mods *modifiers) string {

	requireStructureID(mods)

	return buildKey(GetStructure.String(), strconv.FormatUint(mods.structureID, 10))

}

func structurePathFunc(mods *modifiers) string {

	requireStructureID(mods)

	return fmt.Sprintf(endpoints[GetStructure].Path, mods.structureID)

}
