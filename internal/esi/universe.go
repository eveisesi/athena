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
	GetAncestries(ctx context.Context) ([]*athena.Ancestry, *athena.Etag, *http.Response, error)
	GetBloodlines(ctx context.Context) ([]*athena.Bloodline, *athena.Etag, *http.Response, error)
	GetCategories(ctx context.Context) ([]uint, *athena.Etag, *http.Response, error)
	GetCategory(ctx context.Context, categoryID uint) (*athena.Category, *athena.Etag, *http.Response, error)
	GetConstellation(ctx context.Context, constellationID uint) (*athena.Constellation, *athena.Etag, *http.Response, error)
	GetFactions(ctx context.Context) ([]*athena.Faction, *athena.Etag, *http.Response, error)
	GetGroup(ctx context.Context, groupID uint) (*athena.Group, *athena.Etag, *http.Response, error)
	GetRegions(ctx context.Context) ([]uint, *athena.Etag, *http.Response, error)
	GetRegion(ctx context.Context, regionID uint) (*athena.Region, *athena.Etag, *http.Response, error)
	GetRaces(ctx context.Context) ([]*athena.Race, *athena.Etag, *http.Response, error)
	GetSolarSystem(ctx context.Context, solarSystemID uint) (*athena.SolarSystem, *athena.Etag, *http.Response, error)
	GetStation(ctx context.Context, stationID uint) (*athena.Station, *athena.Etag, *http.Response, error)
	GetStructure(ctx context.Context, structureID uint64, token string) (*athena.Structure, *athena.Etag, *http.Response, error)
	GetType(ctx context.Context, typeID uint) (*athena.Type, *athena.Etag, *http.Response, error)

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
		return nil, nil, fmt.Errorf("failed to marshal ids to byte array: %w", err)
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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, res, fmt.Errorf("post universe name failed with status code %d", res.StatusCode)
	}

	var names = make([]*PostUniverseNamesOK, 0)
	err = json.Unmarshal(b, &names)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return names, res, nil

}

func (s *service) GetAncestries(ctx context.Context) ([]*athena.Ancestry, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetAncestries]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch ancestries, received status code of %d", res.StatusCode)
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

	var ancestries = make([]*athena.Ancestry, 0)
	err = json.Unmarshal(b, &ancestries)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return ancestries, etag, res, nil

}

func ancestriesKeyFunc(mods *modifiers) string {
	return buildKey(GetAncestries.String())
}

func ancestriesPathFunc(mods *modifiers) string {
	return endpoints[GetAncestries].Path
}

func (s *service) GetBloodlines(ctx context.Context) ([]*athena.Bloodline, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetBloodlines]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch bloodlines, received status code of %d", res.StatusCode)
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

	var bloodlines = make([]*athena.Bloodline, 0)
	err = json.Unmarshal(b, &bloodlines)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return bloodlines, etag, res, nil

}

func bloodlinesKeyFunc(mods *modifiers) string {
	return buildKey(GetBloodlines.String())
}

func bloodlinesPathFunc(mods *modifiers) string {
	return endpoints[GetBloodlines].Path
}

func (s *service) GetRaces(ctx context.Context) ([]*athena.Race, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetRaces]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch races, received status code of %d", res.StatusCode)
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

	var races = make([]*athena.Race, 0)
	err = json.Unmarshal(b, &races)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return races, etag, res, nil

}

func racesKeyFunc(mods *modifiers) string {
	return buildKey(GetRaces.String())
}

func racesPathFunc(mods *modifiers) string {
	return endpoints[GetRaces].Path
}

func (s *service) GetFactions(ctx context.Context) ([]*athena.Faction, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetFactions]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch factions, received status code of %d", res.StatusCode)
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

	var factions = make([]*athena.Faction, 0)
	err = json.Unmarshal(b, &factions)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return factions, etag, res, nil

}

func factionsKeyFunc(mods *modifiers) string {
	return buildKey(GetFactions.String())
}

func factionsPathFunc(mods *modifiers) string {
	return endpoints[GetFactions].Path
}

func (s *service) GetCategories(ctx context.Context) ([]uint, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCategories]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch categories, received status code of %d", res.StatusCode)
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

	var categories = make([]uint, 0)
	err = json.Unmarshal(b, &categories)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return categories, etag, res, nil

}

func categoriesKeyFunc(mods *modifiers) string {
	return buildKey(GetCategories.String())
}

func categoriesPathFunc(mods *modifiers) string {
	return endpoints[GetCategories].Path
}

func (s *service) GetCategory(ctx context.Context, categoryID uint) (*athena.Category, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCategory]

	mods := s.modifiers(ModWithCategoryID(categoryID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch category %d, received status code of %d", categoryID, res.StatusCode)
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

	var category = new(athena.Category)
	err = json.Unmarshal(b, category)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	category.ID = categoryID

	return category, etag, res, nil

}

func categoryKeyFunc(mods *modifiers) string {

	// requireCategoryID(mods)

	return buildKey(GetCategory.String(), strconv.FormatUint(uint64(mods.categoryID), 10))

}

func categoryPathFunc(mods *modifiers) string {

	// requireCategoryID(mods)

	return fmt.Sprintf(endpoints[GetCategory].Path, mods.categoryID)

}

func (s *service) GetGroup(ctx context.Context, groupID uint) (*athena.Group, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetGroup]

	mods := s.modifiers(ModWithGroupID(groupID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch group %d, received status code of %d", groupID, res.StatusCode)
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

	var group = new(athena.Group)
	err = json.Unmarshal(b, group)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	group.ID = groupID

	return group, etag, res, nil

}

func groupKeyFunc(mods *modifiers) string {

	// requireGroupID(mods)

	return buildKey(GetGroup.String(), strconv.FormatUint(uint64(mods.groupID), 10))

}

func groupPathFunc(mods *modifiers) string {

	// requireGroupID(mods)

	return fmt.Sprintf(endpoints[GetGroup].Path, mods.groupID)

}

func (s *service) GetType(ctx context.Context, typeID uint) (*athena.Type, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetType]

	mods := s.modifiers(ModWithItemID(typeID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch type %d, received status code of %d", typeID, res.StatusCode)
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

	var item = new(athena.Type)
	err = json.Unmarshal(b, item)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	item.ID = typeID

	return item, etag, res, nil

}

func typeKeyFunc(mods *modifiers) string {

	// requireItemID(mods)

	return buildKey(GetType.String(), strconv.FormatUint(uint64(mods.itemID), 10))

}

func typePathFunc(mods *modifiers) string {

	// requireItemID(mods)

	return fmt.Sprintf(endpoints[GetType].Path, mods.itemID)

}

func (s *service) GetRegions(ctx context.Context) ([]uint, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetRegions]

	mods := s.modifiers()

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch regions, received status code of %d", res.StatusCode)
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

	var ids = make([]uint, 0)
	err = json.Unmarshal(b, &ids)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return ids, etag, res, nil

}

func regionsKeyFunc(mods *modifiers) string {
	return buildKey(GetRegions.String())
}

func regionsPathFunc(mods *modifiers) string {
	return endpoints[GetRegions].Path
}

func (s *service) GetRegion(ctx context.Context, regionID uint) (*athena.Region, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetRegion]

	mods := s.modifiers(ModWithRegionID(regionID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch region %d, received status code of %d", regionID, res.StatusCode)
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

	var region = new(athena.Region)
	err = json.Unmarshal(b, region)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	region.ID = regionID

	return region, etag, res, nil

}

func regionKeyFunc(mods *modifiers) string {

	requireRegionID(mods)

	return buildKey(GetRegion.String(), strconv.FormatUint(uint64(mods.regionID), 10))

}

func regionPathFunc(mods *modifiers) string {

	requireRegionID(mods)

	return fmt.Sprintf(endpoints[GetRegion].Path, mods.regionID)

}

func (s *service) GetConstellation(ctx context.Context, constellationID uint) (*athena.Constellation, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetConstellation]

	mods := s.modifiers(ModWithConstellationID(constellationID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch constellation %d, received status code of %d", constellationID, res.StatusCode)
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

	var constellation = new(athena.Constellation)
	err = json.Unmarshal(b, constellation)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	constellation.ID = constellationID

	return constellation, etag, res, nil

}

func constellationKeyFunc(mods *modifiers) string {

	requireConstellationID(mods)

	return buildKey(GetConstellation.String(), strconv.FormatUint(uint64(mods.constellationID), 10))

}

func constellationPathFunc(mods *modifiers) string {

	requireConstellationID(mods)

	return fmt.Sprintf(endpoints[GetConstellation].Path, mods.constellationID)

}

func (s *service) GetSolarSystem(ctx context.Context, systemID uint) (*athena.SolarSystem, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetSolarSystem]

	mods := s.modifiers(ModWithSystemID(systemID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch solar system %d, received status code of %d", systemID, res.StatusCode)
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

	var system = new(athena.SolarSystem)
	err = json.Unmarshal(b, system)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	system.ID = systemID

	return system, etag, res, nil

}

func solarSystemKeyFunc(mods *modifiers) string {

	requireSystemID(mods)

	return buildKey(GetSolarSystem.String(), strconv.FormatUint(uint64(mods.solarSystemID), 10))

}

func solarSystemPathFunc(mods *modifiers) string {

	requireSystemID(mods)

	return fmt.Sprintf(endpoints[GetSolarSystem].Path, mods.solarSystemID)

}

func (s *service) GetStation(ctx context.Context, stationID uint) (*athena.Station, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetStation]

	mods := s.modifiers(ModWithStationID(stationID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

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
		return nil, etag, res, fmt.Errorf("failed to fetch station %d, received status code of %d", stationID, res.StatusCode)
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

	var station = new(athena.Station)
	err = json.Unmarshal(b, station)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	station.ID = stationID

	return station, etag, res, nil

}

func stationKeyFunc(mods *modifiers) string {

	requireStationID(mods)

	return buildKey(GetStation.String(), strconv.FormatUint(uint64(mods.stationID), 10))

}

func stationPathFunc(mods *modifiers) string {

	requireStationID(mods)

	return fmt.Sprintf(endpoints[GetStation].Path, mods.stationID)

}

func (s *service) GetStructure(ctx context.Context, structureID uint64, token string) (*athena.Structure, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetStructure]

	mods := s.modifiers(ModWithStructureID(structureID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, etag, res, fmt.Errorf("failed to fetch structure %d, received status code of %d", structureID, res.StatusCode)
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

	var structure = new(athena.Structure)
	err = json.Unmarshal(b, structure)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	structure.ID = structureID

	return structure, etag, res, nil

}

func structureKeyFunc(mods *modifiers) string {

	requireStructureID(mods)

	return buildKey(GetStructure.String(), strconv.FormatUint(mods.structureID, 10))

}

func structurePathFunc(mods *modifiers) string {

	requireStructureID(mods)

	return fmt.Sprintf(endpoints[GetStructure].Path, mods.structureID)

}
