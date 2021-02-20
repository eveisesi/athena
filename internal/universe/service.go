package universe

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/sirupsen/logrus"
)

type Service interface {
	InitializeUniverse(options ...OptionFunc) error

	Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error)
	Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error)
	Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error)
	Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error)
	Category(ctx context.Context, id uint) (*athena.Category, error)
	Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error)
	Constellation(ctx context.Context, id uint) (*athena.Constellation, error)
	Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error)
	Faction(ctx context.Context, id uint) (*athena.Faction, error)
	Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error)
	Group(ctx context.Context, id uint) (*athena.Group, error)
	Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error)
	Race(ctx context.Context, id uint) (*athena.Race, error)
	Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error)
	Region(ctx context.Context, id uint) (*athena.Region, error)
	Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error)
	SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error)
	SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error)
	Station(ctx context.Context, id uint) (*athena.Station, error)
	Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error)
	Structure(ctx context.Context, member *athena.Member, id uint64) (*athena.Structure, error)
	Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error)
	Type(ctx context.Context, id uint) (*athena.Type, error)
	Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error)
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	universe athena.UniverseRepository
}

const (
	serviceIdentifier = "Universe Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe athena.UniverseRepository) Service {

	logger.SetFormatter(&logrus.TextFormatter{})

	return &service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		universe: universe,
	}

}

func (s *service) Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error) {

	ancestry, err := s.cache.Ancestry(ctx, id)
	if err != nil {
		return nil, err
	}

	if ancestry != nil {
		return ancestry, nil
	}

	ancestry, err = s.universe.Ancestry(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetAncestry(ctx, ancestry)

	return ancestry, err

}

func (s *service) Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Ancestries",
	})

	ancestries, err := s.cache.Ancestries(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch ancestries from cache")
		return nil, fmt.Errorf("failed to fetch ancestries from cache")
	}

	if len(ancestries) > 0 {
		return ancestries, nil
	}

	ancestries, err = s.universe.Ancestries(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch ancestries from db")
		return nil, fmt.Errorf("failed to fetch ancestries from db")
	}

	err = s.cache.SetAncestries(ctx, ancestries, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache ancestries")
		return nil, fmt.Errorf("failed to cache ancestries")
	}

	return ancestries, nil

}

func (s *service) Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error) {

	bloodline, err := s.cache.Bloodline(ctx, id)
	if err != nil {
		return nil, err
	}

	if bloodline != nil {
		return bloodline, nil
	}

	bloodline, err = s.universe.Bloodline(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetBloodline(ctx, bloodline)

	return bloodline, err

}

func (s *service) Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Bloodlines",
	})

	bloodlines, err := s.cache.Bloodlines(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch bloodlines from cache")
		return nil, fmt.Errorf("failed to fetch bloodlines from cache")
	}

	if len(bloodlines) > 0 {
		return bloodlines, nil
	}

	bloodlines, err = s.universe.Bloodlines(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch bloodlines from db")
		return nil, fmt.Errorf("failed to fetch bloodlines from db")
	}

	err = s.cache.SetBloodlines(ctx, bloodlines, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache bloodlines")
		return nil, fmt.Errorf("failed to cache bloodlines")
	}

	return bloodlines, nil

}

func (s *service) Category(ctx context.Context, id uint) (*athena.Category, error) {

	category, err := s.cache.Category(ctx, id)
	if err != nil {
		return nil, err
	}

	if category != nil {
		return category, nil
	}

	category, err = s.universe.Category(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if category == nil || errors.Is(err, sql.ErrNoRows) {
		category, _, _, err = s.esi.GetCategory(ctx, id)
		if err != nil {
			return nil, err
		}

		category, err = s.universe.CreateCategory(ctx, category)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetCategory(ctx, category)

	return category, err

}

func (s *service) Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Categories",
	})

	categories, err := s.cache.Categories(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch categories from cache")
		return nil, fmt.Errorf("failed to fetch categories from cache")
	}

	if len(categories) > 0 {
		return categories, nil
	}

	categories, err = s.universe.Categories(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch categories from db")
		return nil, fmt.Errorf("failed to fetch categories from db")
	}

	err = s.cache.SetCategories(ctx, categories, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache categories")
		return nil, fmt.Errorf("failed to cache categories")
	}

	return categories, nil

}

func (s *service) Constellation(ctx context.Context, id uint) (*athena.Constellation, error) {

	constellation, err := s.cache.Constellation(ctx, id)
	if err != nil {
		return nil, err
	}

	if constellation != nil {
		return constellation, nil
	}

	constellation, err = s.universe.Constellation(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if constellation == nil || errors.Is(err, sql.ErrNoRows) {
		constellation, _, _, err = s.esi.GetConstellation(ctx, id)
		if err != nil {
			return nil, err
		}

		constellation, err = s.universe.CreateConstellation(ctx, constellation)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetConstellation(ctx, constellation)

	return constellation, err

}

func (s *service) Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Constellations",
	})

	constellations, err := s.cache.Constellations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch constellations from cache")
		return nil, fmt.Errorf("failed to fetch constellations from cache")
	}

	if len(constellations) > 0 {
		return constellations, nil
	}

	constellations, err = s.universe.Constellations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch constellations from db")
		return nil, fmt.Errorf("failed to fetch constellations from db")
	}

	err = s.cache.SetConstellations(ctx, constellations, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache constellations")
		return nil, fmt.Errorf("failed to cache constellations")
	}

	return constellations, nil

}

func (s *service) Faction(ctx context.Context, id uint) (*athena.Faction, error) {
	faction, err := s.cache.Faction(ctx, id)
	if err != nil {
		return nil, err
	}

	if faction != nil {
		return faction, nil
	}

	faction, err = s.universe.Faction(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetFaction(ctx, faction)

	return faction, err
}

func (s *service) Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Factions",
	})

	factions, err := s.cache.Factions(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch factions from cache")
		return nil, fmt.Errorf("failed to fetch factions from cache")
	}

	if len(factions) > 0 {
		return factions, nil
	}

	factions, err = s.universe.Factions(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch factions from db")
		return nil, fmt.Errorf("failed to fetch factions from db")
	}

	err = s.cache.SetFactions(ctx, factions, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache factions")
		return nil, fmt.Errorf("failed to cache factions")
	}

	return factions, nil

}

func (s *service) Group(ctx context.Context, id uint) (*athena.Group, error) {

	group, err := s.cache.Group(ctx, id)
	if err != nil {
		return nil, err
	}

	if group != nil {
		return group, nil
	}

	group, err = s.universe.Group(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if group == nil || errors.Is(err, sql.ErrNoRows) {
		group, _, _, err = s.esi.GetGroup(ctx, id)
		if err != nil {
			return nil, err
		}

		group, err = s.universe.CreateGroup(ctx, group)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetGroup(ctx, group)

	return group, err

}

func (s *service) Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Groups",
	})

	groups, err := s.cache.Groups(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch groups from cache")
		return nil, fmt.Errorf("failed to fetch groups from cache")
	}

	if len(groups) > 0 {
		return groups, nil
	}

	groups, err = s.universe.Groups(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch groups from db")
		return nil, fmt.Errorf("failed to fetch groups from db")
	}

	err = s.cache.SetGroups(ctx, groups, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache groups")
		return nil, fmt.Errorf("failed to cache groups")
	}

	return groups, nil

}

func (s *service) Race(ctx context.Context, id uint) (*athena.Race, error) {

	race, err := s.cache.Race(ctx, id)
	if err != nil {
		return nil, err
	}

	if race != nil {
		return race, nil
	}

	race, err = s.universe.Race(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetRace(ctx, race)

	return race, err

}

func (s *service) Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Races",
	})

	races, err := s.cache.Races(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch races from cache")
		return nil, fmt.Errorf("failed to fetch races from cache")
	}

	if len(races) > 0 {
		return races, nil
	}

	races, err = s.universe.Races(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch races from db")
		return nil, fmt.Errorf("failed to fetch races from db")
	}

	err = s.cache.SetRaces(ctx, races, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache races")
		return nil, fmt.Errorf("failed to cache races")
	}

	return races, nil

}

func (s *service) Region(ctx context.Context, id uint) (*athena.Region, error) {

	region, err := s.cache.Region(ctx, id)
	if err != nil {
		return nil, err
	}

	if region != nil {
		return region, nil
	}

	region, err = s.universe.Region(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if region == nil || errors.Is(err, sql.ErrNoRows) {
		region, _, _, err = s.esi.GetRegion(ctx, id)
		if err != nil {
			return nil, err
		}

		region, err = s.universe.CreateRegion(ctx, region)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetRegion(ctx, region)

	return region, err

}

func (s *service) Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Regions",
	})

	regions, err := s.cache.Regions(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch regions from cache")
		return nil, fmt.Errorf("failed to fetch regions from cache")
	}

	if len(regions) > 0 {
		return regions, nil
	}

	regions, err = s.universe.Regions(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch regions from db")
		return nil, fmt.Errorf("failed to fetch regions from db")
	}

	err = s.cache.SetRegions(ctx, regions, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache regions")
		return nil, fmt.Errorf("failed to cache regions")
	}

	return regions, nil

}

func (s *service) SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error) {

	solarSystem, err := s.cache.SolarSystem(ctx, id)
	if err != nil {
		return nil, err
	}

	if solarSystem != nil {
		return solarSystem, nil
	}

	solarSystem, err = s.universe.SolarSystem(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if solarSystem == nil || errors.Is(err, sql.ErrNoRows) {
		solarSystem, _, _, err = s.esi.GetSolarSystem(ctx, id)
		if err != nil {
			return nil, err
		}

		solarSystem, err = s.universe.CreateSolarSystem(ctx, solarSystem)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetSolarSystem(ctx, solarSystem)

	return solarSystem, err

}

func (s *service) SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "SolarSystems",
	})

	systems, err := s.cache.SolarSystems(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch systems from cache")
		return nil, fmt.Errorf("failed to fetch systems from cache")
	}

	if len(systems) > 0 {
		return systems, nil
	}

	systems, err = s.universe.SolarSystems(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch systems from db")
		return nil, fmt.Errorf("failed to fetch systems from db")
	}

	err = s.cache.SetSolarSystems(ctx, systems, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache systems")
		return nil, fmt.Errorf("failed to cache systems")
	}

	return systems, nil

}

func (s *service) Station(ctx context.Context, id uint) (*athena.Station, error) {

	station, err := s.cache.Station(ctx, id)
	if err != nil {
		return nil, err
	}

	if station != nil {
		return station, nil
	}

	station, err = s.universe.Station(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if station == nil || errors.Is(err, sql.ErrNoRows) {
		station, _, _, err = s.esi.GetStation(ctx, id)
		if err != nil {
			return nil, err
		}

		station, err = s.universe.CreateStation(ctx, station)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetStation(ctx, station)

	return station, err

}

func (s *service) Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Stations",
	})

	stations, err := s.cache.Stations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch stations from cache")
		return nil, fmt.Errorf("failed to fetch stations from cache")
	}

	if len(stations) > 0 {
		return stations, nil
	}

	stations, err = s.universe.Stations(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch stations from db")
		return nil, fmt.Errorf("failed to fetch stations from db")
	}

	err = s.cache.SetStations(ctx, stations, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache stations")
		return nil, fmt.Errorf("failed to cache stations")
	}

	return stations, nil

}

func (s *service) Structure(ctx context.Context, member *athena.Member, id uint64) (*athena.Structure, error) {

	structure, err := s.cache.Structure(ctx, id)
	if err != nil {
		return nil, err
	}

	if structure != nil {
		return structure, nil
	}

	structure, err = s.universe.Structure(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if structure == nil || errors.Is(err, sql.ErrNoRows) {
		// TODO: Deliver a Concreate Error from ESI Package and insert th is into
		structure, _, _, err = s.esi.GetStructure(ctx, id, member.AccessToken.String)
		if err != nil {
			return nil, err
		}

		structure, err = s.universe.CreateStructure(ctx, structure)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetStructure(ctx, structure)

	return structure, err

}

func (s *service) Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Structures",
	})

	structures, err := s.cache.Structures(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch structures from cache")
		return nil, fmt.Errorf("failed to fetch structures from cache")
	}

	if len(structures) > 0 {
		return structures, nil
	}

	structures, err = s.universe.Structures(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch structures from db")
		return nil, fmt.Errorf("failed to fetch structures from db")
	}

	err = s.cache.SetStructures(ctx, structures, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache structures")
		return nil, fmt.Errorf("failed to cache structures")
	}

	return structures, nil

}

func (s *service) Type(ctx context.Context, id uint) (*athena.Type, error) {

	item, err := s.cache.Type(ctx, id)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return item, nil
	}

	item, err = s.universe.Type(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if item == nil || errors.Is(err, sql.ErrNoRows) {
		item, _, _, err = s.esi.GetType(ctx, id)
		if err != nil {
			return nil, err
		}

		item, err = s.universe.CreateType(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetType(ctx, item)

	return item, err

}

func (s *service) Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "Types",
	})

	types, err := s.cache.Types(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch types from cache")
		return nil, fmt.Errorf("failed to fetch types from cache")
	}

	if len(types) > 0 {
		return types, nil
	}

	types, err = s.universe.Types(ctx, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to fetch types from db")
		return nil, fmt.Errorf("failed to fetch types from db")
	}

	err = s.cache.SetTypes(ctx, types, operators...)
	if err != nil {
		entry.WithError(err).Error("failed to cache types")
		return nil, fmt.Errorf("failed to cache types")
	}

	return types, nil

}
