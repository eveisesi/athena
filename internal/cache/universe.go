package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type universeService interface {
	Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error)
	SetAncestry(ctx context.Context, ancestry *athena.Ancestry, optionFuncs ...OptionFunc) error
	Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error)
	SetBloodline(ctx context.Context, bloodline *athena.Bloodline, optionFuncs ...OptionFunc) error
	Race(ctx context.Context, id uint) (*athena.Race, error)
	SetRace(ctx context.Context, race *athena.Race, optionFuncs ...OptionFunc) error
	Faction(ctx context.Context, id uint) (*athena.Faction, error)
	SetFaction(ctx context.Context, faction *athena.Faction, optionFuncs ...OptionFunc) error
	Region(ctx context.Context, id uint) (*athena.Region, error)
	SetRegion(ctx context.Context, region *athena.Region, optionFuncs ...OptionFunc) error
	Constellation(ctx context.Context, id uint) (*athena.Constellation, error)
	SetConstellation(ctx context.Context, constellation *athena.Constellation, optionFuncs ...OptionFunc) error
	SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error)
	SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem, optionFuncs ...OptionFunc) error
	AsteroidBelt(ctx context.Context, id uint) (*athena.AsteroidBelt, error)
	SetAsteroidBelt(ctx context.Context, belt *athena.AsteroidBelt, optionFuncs ...OptionFunc) error
	Moon(ctx context.Context, id uint) (*athena.Moon, error)
	SetMoon(ctx context.Context, moon *athena.Moon, optionFuncs ...OptionFunc) error
	Planet(ctx context.Context, id uint) (*athena.Planet, error)
	SetPlanet(ctx context.Context, planet *athena.Planet, optionFuncs ...OptionFunc) error
	Station(ctx context.Context, id uint) (*athena.Station, error)
	SetStation(ctx context.Context, station *athena.Station, optionFuncs ...OptionFunc) error
	Structure(ctx context.Context, id uint64) (*athena.Structure, error)
	SetStructure(ctx context.Context, structure *athena.Structure, optionFuncs ...OptionFunc) error
	Category(ctx context.Context, id uint) (*athena.Category, error)
	SetCategory(ctx context.Context, category *athena.Category, optionFuncs ...OptionFunc) error
	Group(ctx context.Context, id uint) (*athena.Group, error)
	SetGroup(ctx context.Context, group *athena.Group, optionFuncs ...OptionFunc) error
	Type(ctx context.Context, id uint) (*athena.Type, error)
	SetType(ctx context.Context, item *athena.Type, optionFuncs ...OptionFunc) error
}

const (
	keyAncestry      = "athena::ancestry::%d"
	keyAsteroidBelt  = "athena::anstroidBelt::%d"
	keyBloodline     = "athena::bloodline::%d"
	keyRace          = "athena::race::%d"
	keyFaction       = "athena::faction::%d"
	keyCategory      = "athena::category::%d"
	keyGroup         = "athena::group::%d"
	keyMoon          = "athena::moon::%d"
	keyType          = "athena::type::%d"
	keyRegion        = "athena::region::%d"
	keyConstellation = "athena::constellation::%d"
	keyPlanet        = "athena::planet::%d"
	keySolarSystem   = "athena::solar_system::%d"
	keyStation       = "athena::station::%d"
	keyStructure     = "athena::structure::%d"
)

func (s *service) Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error) {

	key := fmt.Sprintf(keyAncestry, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var ancestry = new(athena.Ancestry)
	err = json.Unmarshal(result, &ancestry)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return ancestry, nil

}

func (s *service) SetAncestry(ctx context.Context, ancestry *athena.Ancestry, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyAncestry, ancestry.ID)
	data, err := json.Marshal(ancestry)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error) {

	key := fmt.Sprintf(keyBloodline, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var bloodline = new(athena.Bloodline)
	err = json.Unmarshal(result, &bloodline)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return bloodline, nil

}

func (s *service) SetBloodline(ctx context.Context, bloodline *athena.Bloodline, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyBloodline, bloodline.ID)
	data, err := json.Marshal(bloodline)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Race(ctx context.Context, id uint) (*athena.Race, error) {

	key := fmt.Sprintf(keyRace, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var race = new(athena.Race)
	err = json.Unmarshal(result, &race)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return race, nil

}

func (s *service) SetRace(ctx context.Context, race *athena.Race, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyRace, race.ID)
	data, err := json.Marshal(race)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Faction(ctx context.Context, id uint) (*athena.Faction, error) {

	key := fmt.Sprintf(keyFaction, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var faction = new(athena.Faction)
	err = json.Unmarshal(result, &faction)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return faction, nil

}

func (s *service) SetFaction(ctx context.Context, faction *athena.Faction, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyFaction, faction.ID)
	data, err := json.Marshal(faction)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Region(ctx context.Context, id uint) (*athena.Region, error) {

	key := fmt.Sprintf(keyRegion, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var region = new(athena.Region)
	err = json.Unmarshal(result, &region)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return region, nil

}

func (s *service) SetRegion(ctx context.Context, region *athena.Region, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyRegion, region.ID)
	data, err := json.Marshal(region)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Constellation(ctx context.Context, id uint) (*athena.Constellation, error) {

	key := fmt.Sprintf(keyConstellation, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var constellation = new(athena.Constellation)

	err = json.Unmarshal(result, &constellation)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return constellation, nil

}

func (s *service) SetConstellation(ctx context.Context, constellation *athena.Constellation, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyConstellation, constellation.ID)
	data, err := json.Marshal(constellation)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error) {

	key := fmt.Sprintf(keySolarSystem, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var solarSystem = new(athena.SolarSystem)

	err = json.Unmarshal(result, &solarSystem)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return solarSystem, nil

}

func (s *service) SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keySolarSystem, solarSystem.ID)
	data, err := json.Marshal(solarSystem)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Planet(ctx context.Context, id uint) (*athena.Planet, error) {

	key := fmt.Sprintf(keyPlanet, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer Fialed to fetch results from cache for key %s: %w]", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var planet = new(athena.Planet)
	err = json.Unmarshal(result, planet)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarsahl results for key %s on struct: %w", key, err)
	}

	return planet, nil

}

func (s *service) SetPlanet(ctx context.Context, planet *athena.Planet, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyPlanet, planet.ID)
	data, err := json.Marshal(planet)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Moon(ctx context.Context, id uint) (*athena.Moon, error) {

	key := fmt.Sprintf(keyMoon, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer Fialed to fetch results from cache for key %s: %w]", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var moon = new(athena.Moon)
	err = json.Unmarshal(result, moon)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarsahl results for key %s on struct: %w", key, err)
	}

	return moon, nil

}

func (s *service) SetMoon(ctx context.Context, moon *athena.Moon, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyMoon, moon.ID)
	data, err := json.Marshal(moon)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) AsteroidBelt(ctx context.Context, id uint) (*athena.AsteroidBelt, error) {

	key := fmt.Sprintf(keyAsteroidBelt, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer Fialed to fetch results from cache for key %s: %w]", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var belt = new(athena.AsteroidBelt)
	err = json.Unmarshal(result, belt)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarsahl results for key %s on struct: %w", key, err)
	}

	return belt, nil

}

func (s *service) SetAsteroidBelt(ctx context.Context, belt *athena.AsteroidBelt, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyAsteroidBelt, belt.ID)
	data, err := json.Marshal(belt)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Station(ctx context.Context, id uint) (*athena.Station, error) {

	key := fmt.Sprintf(keyStation, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var station = new(athena.Station)

	err = json.Unmarshal(result, &station)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return station, nil

}

func (s *service) SetStation(ctx context.Context, station *athena.Station, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyStation, station.ID)
	data, err := json.Marshal(station)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyStation, station.ID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Structure(ctx context.Context, id uint64) (*athena.Structure, error) {

	key := fmt.Sprintf(keyStructure, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var structure = new(athena.Structure)

	err = json.Unmarshal(result, &structure)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return structure, nil

}

func (s *service) SetStructure(ctx context.Context, structure *athena.Structure, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyStructure, structure.ID)
	data, err := json.Marshal(structure)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Category(ctx context.Context, id uint) (*athena.Category, error) {

	key := fmt.Sprintf(keyCategory, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var category = new(athena.Category)

	err = json.Unmarshal(result, &category)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return category, nil

}

func (s *service) SetCategory(ctx context.Context, category *athena.Category, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyCategory, category.ID)
	data, err := json.Marshal(category)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Group(ctx context.Context, id uint) (*athena.Group, error) {

	key := fmt.Sprintf(keyGroup, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var group = new(athena.Group)

	err = json.Unmarshal(result, &group)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal results for key %s on struct: %w", key, err)
	}

	return group, nil

}

func (s *service) SetGroup(ctx context.Context, group *athena.Group, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyGroup, group.ID)
	data, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Type(ctx context.Context, id uint) (*athena.Type, error) {

	key := fmt.Sprintf(keyType, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var item = new(athena.Type)
	err = json.Unmarshal(result, &item)
	if err != nil {
		return nil, err
	}

	return item, nil

}

func (s *service) SetType(ctx context.Context, item *athena.Type, optionFuncs ...OptionFunc) error {

	key := fmt.Sprintf(keyType, item.ID)
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}
