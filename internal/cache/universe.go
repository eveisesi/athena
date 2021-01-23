package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type universeService interface {
	Ancestry(ctx context.Context, id int) (*athena.Ancestry, error)
	SetAncestry(ctx context.Context, ancestry *athena.Ancestry, optionFuncs ...OptionFunc) error
	Bloodline(ctx context.Context, id int) (*athena.Bloodline, error)
	SetBloodline(ctx context.Context, bloodline *athena.Bloodline, optionFuncs ...OptionFunc) error
	Race(ctx context.Context, id int) (*athena.Race, error)
	SetRace(ctx context.Context, race *athena.Race, optionFuncs ...OptionFunc) error
	Faction(ctx context.Context, id int) (*athena.Faction, error)
	SetFaction(ctx context.Context, faction *athena.Faction, optionFuncs ...OptionFunc) error
	Region(ctx context.Context, id int) (*athena.Region, error)
	SetRegion(ctx context.Context, region *athena.Region, optionFuncs ...OptionFunc) error
	Constellation(ctx context.Context, id int) (*athena.Constellation, error)
	SetConstellation(ctx context.Context, constellation *athena.Constellation, optionFuncs ...OptionFunc) error
	SolarSystem(ctx context.Context, id int) (*athena.SolarSystem, error)
	SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem, optionFuncs ...OptionFunc) error
	Station(ctx context.Context, id int) (*athena.Station, error)
	SetStation(ctx context.Context, station *athena.Station, optionFuncs ...OptionFunc) error
	Structure(ctx context.Context, id int64) (*athena.Structure, error)
	SetStructure(ctx context.Context, station *athena.Structure, optionFuncs ...OptionFunc) error
	Category(ctx context.Context, id int) (*athena.Category, error)
	SetCategory(ctx context.Context, category *athena.Category, optionFuncs ...OptionFunc) error
	Group(ctx context.Context, id int) (*athena.Group, error)
	SetGroup(ctx context.Context, group *athena.Group, optionFuncs ...OptionFunc) error
	Type(ctx context.Context, id int) (*athena.Type, error)
	SetType(ctx context.Context, item *athena.Type, optionFuncs ...OptionFunc) error
}

const (
	keyAncestry      = "athena::ancestry::%d"
	keyBloodline     = "athena::bloodline::%d"
	keyRace          = "athena::race::%d"
	keyFaction       = "athena::faction::%d"
	keyCategory      = "athena::category::%d"
	keyGroup         = "athena::group::%d"
	keyType          = "athena::type::%d"
	keyRegion        = "athena::region::%d"
	keyConstellation = "athena::constellation::%d"
	keySolarSystem   = "athena::solar_system::%d"
	keyStation       = "athena::station::%d"
	keyStructure     = "athena::structure::%d"
)

func (s *service) Ancestry(ctx context.Context, id int) (*athena.Ancestry, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyAncestry, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var ancestry = new(athena.Ancestry)

	err = json.Unmarshal([]byte(result), &ancestry)
	if err != nil {
		return nil, err
	}

	return ancestry, nil

}

func (s *service) SetAncestry(ctx context.Context, ancestry *athena.Ancestry, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(ancestry)
	if err != nil {
		return fmt.Errorf("failed to marshal ancestry: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyAncestry, ancestry.AncestryID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Bloodline(ctx context.Context, id int) (*athena.Bloodline, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyBloodline, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var bloodline = new(athena.Bloodline)

	err = json.Unmarshal([]byte(result), &bloodline)
	if err != nil {
		return nil, err
	}

	return bloodline, nil

}

func (s *service) SetBloodline(ctx context.Context, bloodline *athena.Bloodline, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(bloodline)
	if err != nil {
		return fmt.Errorf("failed to marshal bloodline: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyBloodline, bloodline.BloodlineID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Race(ctx context.Context, id int) (*athena.Race, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyRace, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var race = new(athena.Race)

	err = json.Unmarshal([]byte(result), &race)
	if err != nil {
		return nil, err
	}

	return race, nil

}

func (s *service) SetRace(ctx context.Context, race *athena.Race, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(race)
	if err != nil {
		return fmt.Errorf("failed to marshal race: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyRace, race.RaceID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Faction(ctx context.Context, id int) (*athena.Faction, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyFaction, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var faction = new(athena.Faction)

	err = json.Unmarshal([]byte(result), &faction)
	if err != nil {
		return nil, err
	}

	return faction, nil

}

func (s *service) SetFaction(ctx context.Context, faction *athena.Faction, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(faction)
	if err != nil {
		return fmt.Errorf("failed to marshal faction: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyFaction, faction.FactionID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Region(ctx context.Context, id int) (*athena.Region, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyRegion, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var region = new(athena.Region)

	err = json.Unmarshal([]byte(result), &region)
	if err != nil {
		return nil, err
	}

	return region, nil

}

func (s *service) SetRegion(ctx context.Context, region *athena.Region, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(region)
	if err != nil {
		return fmt.Errorf("failed to marshal region: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyRegion, region.RegionID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Constellation(ctx context.Context, id int) (*athena.Constellation, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyConstellation, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var constellation = new(athena.Constellation)

	err = json.Unmarshal([]byte(result), &constellation)
	if err != nil {
		return nil, err
	}

	return constellation, nil

}

func (s *service) SetConstellation(ctx context.Context, constellation *athena.Constellation, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(constellation)
	if err != nil {
		return fmt.Errorf("failed to marshal constellation: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyConstellation, constellation.ConstellationID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) SolarSystem(ctx context.Context, id int) (*athena.SolarSystem, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keySolarSystem, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var solarSystem = new(athena.SolarSystem)

	err = json.Unmarshal([]byte(result), &solarSystem)
	if err != nil {
		return nil, err
	}

	return solarSystem, nil

}

func (s *service) SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(solarSystem)
	if err != nil {
		return fmt.Errorf("failed to marshal solarSystem: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keySolarSystem, solarSystem.SystemID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) Station(ctx context.Context, id int) (*athena.Station, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyStation, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var station = new(athena.Station)

	err = json.Unmarshal([]byte(result), &station)
	if err != nil {
		return nil, err
	}

	return station, nil

}

func (s *service) SetStation(ctx context.Context, station *athena.Station, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(station)
	if err != nil {
		return fmt.Errorf("failed to marshal station: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyStation, station.StationID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) Structure(ctx context.Context, id int64) (*athena.Structure, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyStructure, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var structure = new(athena.Structure)

	err = json.Unmarshal([]byte(result), &structure)
	if err != nil {
		return nil, err
	}

	return structure, nil

}

func (s *service) SetStructure(ctx context.Context, structure *athena.Structure, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(structure)
	if err != nil {
		return fmt.Errorf("failed to marshal structure: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyStructure, structure.StructureID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) Category(ctx context.Context, id int) (*athena.Category, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCategory, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var category = new(athena.Category)

	err = json.Unmarshal([]byte(result), &category)
	if err != nil {
		return nil, err
	}

	return category, nil

}

func (s *service) SetCategory(ctx context.Context, category *athena.Category, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(category)
	if err != nil {
		return fmt.Errorf("failed to marshal category: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCategory, category.CategoryID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Group(ctx context.Context, id int) (*athena.Group, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyGroup, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var group = new(athena.Group)

	err = json.Unmarshal([]byte(result), &group)
	if err != nil {
		return nil, err
	}

	return group, nil

}

func (s *service) SetGroup(ctx context.Context, group *athena.Group, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("failed to marshal group: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyGroup, group.GroupID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) Type(ctx context.Context, id int) (*athena.Type, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyType, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var item = new(athena.Type)

	err = json.Unmarshal([]byte(result), &item)
	if err != nil {
		return nil, err
	}

	return item, nil

}

func (s *service) SetType(ctx context.Context, item *athena.Type, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyType, item.TypeID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}
