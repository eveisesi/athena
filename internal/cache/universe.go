package cache

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type universeService interface {
	Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error)
	SetAncestry(ctx context.Context, ancestry *athena.Ancestry) error
	Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error)
	SetAncestries(ctx context.Context, ancestries []*athena.Ancestry, operators ...*athena.Operator) error
	Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error)
	SetBloodline(ctx context.Context, bloodline *athena.Bloodline) error
	Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error)
	SetBloodlines(ctx context.Context, records []*athena.Bloodline, operators ...*athena.Operator) error
	Race(ctx context.Context, id uint) (*athena.Race, error)
	SetRace(ctx context.Context, race *athena.Race) error
	Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error)
	SetRaces(ctx context.Context, records []*athena.Race, operators ...*athena.Operator) error
	Faction(ctx context.Context, id uint) (*athena.Faction, error)
	SetFaction(ctx context.Context, faction *athena.Faction) error
	Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error)
	SetFactions(ctx context.Context, records []*athena.Faction, operators ...*athena.Operator) error
	Region(ctx context.Context, id uint) (*athena.Region, error)
	SetRegion(ctx context.Context, region *athena.Region) error
	Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error)
	SetRegions(ctx context.Context, records []*athena.Region, operators ...*athena.Operator) error
	Constellation(ctx context.Context, id uint) (*athena.Constellation, error)
	SetConstellation(ctx context.Context, constellation *athena.Constellation) error
	Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error)
	SetConstellations(ctx context.Context, records []*athena.Constellation, operators ...*athena.Operator) error
	SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error)
	SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) error
	SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error)
	SetSolarSystems(ctx context.Context, records []*athena.SolarSystem, operators ...*athena.Operator) error
	Station(ctx context.Context, id uint) (*athena.Station, error)
	SetStation(ctx context.Context, station *athena.Station) error
	Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error)
	SetStations(ctx context.Context, records []*athena.Station, operators ...*athena.Operator) error
	Structure(ctx context.Context, id uint64) (*athena.Structure, error)
	SetStructure(ctx context.Context, structure *athena.Structure) error
	Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error)
	SetStructures(ctx context.Context, records []*athena.Structure, operators ...*athena.Operator) error
	Category(ctx context.Context, id uint) (*athena.Category, error)
	SetCategory(ctx context.Context, category *athena.Category) error
	Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error)
	SetCategories(ctx context.Context, records []*athena.Category, operators ...*athena.Operator) error
	Group(ctx context.Context, id uint) (*athena.Group, error)
	SetGroup(ctx context.Context, group *athena.Group) error
	Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error)
	SetGroups(ctx context.Context, records []*athena.Group, operators ...*athena.Operator) error
	Type(ctx context.Context, id uint) (*athena.Type, error)
	SetType(ctx context.Context, item *athena.Type) error
	Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error)
	SetTypes(ctx context.Context, records []*athena.Type, operators ...*athena.Operator) error
}

const (
	keyAncestry       = "athena::ancestry::%d"
	keyAncestries     = "athena::ancestries::%x"
	keyBloodline      = "athena::bloodline::%d"
	keyBloodlines     = "athena::bloodlines::%x"
	keyRace           = "athena::race::%d"
	keyRaces          = "athena::races::%x"
	keyFaction        = "athena::faction::%d"
	keyFactions       = "athena::factions::%x"
	keyCategory       = "athena::category::%d"
	keyCategories     = "athena::categories::%x"
	keyGroup          = "athena::group::%d"
	keyGroups         = "athena::groups::%x"
	keyType           = "athena::type::%d"
	keyTypes          = "athena::types::%x"
	keyRegion         = "athena::region::%d"
	keyRegions        = "athena::regions::%x"
	keyConstellation  = "athena::constellation::%d"
	keyConstellations = "athena::constellations::%x"
	keySolarSystem    = "athena::solar_system::%d"
	keySolarSystems   = "athena::solar_systems::%x"
	keyStation        = "athena::station::%d"
	keyStations       = "athena::stations::%x"
	keyStructure      = "athena::structure::%d"
	keyStructures     = "athena::structures::%x"
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

func (s *service) SetAncestry(ctx context.Context, ancestry *athena.Ancestry) error {

	key := fmt.Sprintf(keyAncestry, ancestry.ID)
	data, err := json.Marshal(ancestry)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyAncestries, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var ancestries = make([]*athena.Ancestry, 0, len(members))
	for _, member := range members {
		var ancestry = new(athena.Ancestry)
		err = json.Unmarshal([]byte(member), ancestry)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal ancestry: %w", err)
		}

		ancestries = append(ancestries, ancestry)
	}

	return ancestries, nil

}

func (s *service) SetAncestries(ctx context.Context, ancestries []*athena.Ancestry, operators ...*athena.Operator) error {

	members := make([]string, 0, len(ancestries))
	for _, ancestry := range ancestries {
		b, err := json.Marshal(ancestry)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal ancestry for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyAncestries, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache labels: %w", err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetBloodline(ctx context.Context, bloodline *athena.Bloodline) error {

	key := fmt.Sprintf(keyBloodline, bloodline.ID)
	data, err := json.Marshal(bloodline)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyBloodlines, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Bloodline, 0, len(members))
	for _, member := range members {
		var result = new(athena.Bloodline)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetBloodlines(ctx context.Context, records []*athena.Bloodline, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyBloodlines, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetRace(ctx context.Context, race *athena.Race) error {

	key := fmt.Sprintf(keyRace, race.ID)
	data, err := json.Marshal(race)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyRaces, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Race, 0, len(members))
	for _, member := range members {
		var result = new(athena.Race)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetRaces(ctx context.Context, records []*athena.Race, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyRaces, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetFaction(ctx context.Context, faction *athena.Faction) error {

	key := fmt.Sprintf(keyFaction, faction.ID)
	data, err := json.Marshal(faction)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyFactions, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Faction, 0, len(members))
	for _, member := range members {
		var result = new(athena.Faction)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetFactions(ctx context.Context, records []*athena.Faction, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyFactions, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetRegion(ctx context.Context, region *athena.Region) error {

	key := fmt.Sprintf(keyRegion, region.ID)
	data, err := json.Marshal(region)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyRegions, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Region, 0, len(members))
	for _, member := range members {
		var result = new(athena.Region)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetRegions(ctx context.Context, records []*athena.Region, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyRegions, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetConstellation(ctx context.Context, constellation *athena.Constellation) error {

	key := fmt.Sprintf(keyConstellation, constellation.ID)
	data, err := json.Marshal(constellation)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyConstellations, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Constellation, 0, len(members))
	for _, member := range members {
		var result = new(athena.Constellation)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetConstellations(ctx context.Context, records []*athena.Constellation, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyConstellations, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) error {

	key := fmt.Sprintf(keySolarSystem, solarSystem.ID)
	data, err := json.Marshal(solarSystem)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keySolarSystems, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.SolarSystem, 0, len(members))
	for _, member := range members {
		var result = new(athena.SolarSystem)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetSolarSystems(ctx context.Context, records []*athena.SolarSystem, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keySolarSystems, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetStation(ctx context.Context, station *athena.Station) error {

	key := fmt.Sprintf(keyStation, station.ID)
	data, err := json.Marshal(station)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyStation, station.ID), data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyStations, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Station, 0, len(members))
	for _, member := range members {
		var result = new(athena.Station)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetStations(ctx context.Context, records []*athena.Station, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyStations, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetStructure(ctx context.Context, structure *athena.Structure) error {

	key := fmt.Sprintf(keyStructure, structure.ID)
	data, err := json.Marshal(structure)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyStructures, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Structure, 0, len(members))
	for _, member := range members {
		var result = new(athena.Structure)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetStructures(ctx context.Context, records []*athena.Structure, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyStructures, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetCategory(ctx context.Context, category *athena.Category) error {

	key := fmt.Sprintf(keyCategory, category.ID)
	data, err := json.Marshal(category)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyCategories, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Category, 0, len(members))
	for _, member := range members {
		var result = new(athena.Category)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetCategories(ctx context.Context, records []*athena.Category, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyCategories, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetGroup(ctx context.Context, group *athena.Group) error {

	key := fmt.Sprintf(keyGroup, group.ID)
	data, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal struct for key %s: %w", key, err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyGroups, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Group, 0, len(members))
	for _, member := range members {
		var result = new(athena.Group)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetGroups(ctx context.Context, records []*athena.Group, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyGroups, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
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

func (s *service) SetType(ctx context.Context, item *athena.Type) error {

	key := fmt.Sprintf(keyType, item.ID)
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = s.client.Set(ctx, key, data, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error) {

	data, err := json.Marshal(operators)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyTypes, sha1.Sum(data))
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	var results = make([]*athena.Type, 0, len(members))
	for _, member := range members {
		var result = new(athena.Type)
		err = json.Unmarshal([]byte(member), result)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal result: %w", err)
		}

		results = append(results, result)
	}

	return results, nil

}

func (s *service) SetTypes(ctx context.Context, records []*athena.Type, operators ...*athena.Operator) error {

	members := make([]string, 0, len(records))
	for _, record := range records {
		b, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal record for cache: %w", err)
		}

		members = append(members, string(b))
	}

	data, err := json.Marshal(operators)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal operators: %w", err)
	}

	key := fmt.Sprintf(keyTypes, sha1.Sum(data))
	_, err = s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}
