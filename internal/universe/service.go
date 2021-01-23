package universe

import (
	"context"
	"fmt"
	"sync"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	InitializeUniverse() error

	Ancestry(ctx context.Context, id int) (*athena.Ancestry, error)
	Bloodline(ctx context.Context, id int) (*athena.Bloodline, error)
	Category(ctx context.Context, id int) (*athena.Category, error)
	Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error)
	Constellation(ctx context.Context, id int) (*athena.Constellation, error)
	Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error)
	Faction(ctx context.Context, id int) (*athena.Faction, error)
	Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error)
	Group(ctx context.Context, id int) (*athena.Group, error)
	Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error)
	Race(ctx context.Context, id int) (*athena.Race, error)
	Region(ctx context.Context, id int) (*athena.Region, error)
	Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error)
	SolarSystem(ctx context.Context, id int) (*athena.SolarSystem, error)
	SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error)
	Station(ctx context.Context, id int) (*athena.Station, error)
	Structure(ctx context.Context, member *athena.Member, id int64) (*athena.Structure, error)
	Type(ctx context.Context, id int) (*athena.Type, error)
	Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error)
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	universe athena.UniverseRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe athena.UniverseRepository) Service {
	// localLogger := *logger
	logger.SetFormatter(&logrus.TextFormatter{})

	return &service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		universe: universe,
	}
}

func (s *service) InitializeUniverse() error {
	var wg = &sync.WaitGroup{}
	var ctx = context.Background()
	// races, _, err := s.esi.GetUniverseRaces(ctx, []*athena.Race{})
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch races from ESI: %w", err)
	// }

	p := mpb.New(mpb.WithWaitGroup(wg))
	// name := "Races"
	// bar := p.AddBar(int64(len(races)),
	// 	mpb.PrependDecorators(
	// 		// simple name decorator
	// 		decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
	// 		// decor.DSyncWidth bit enables column width synchronization
	// 		decor.CountersNoUnit("%d / %d"),
	// 	),
	// 	mpb.AppendDecorators(decor.Percentage()),
	// )

	// for _, race := range races {
	// 	entry := s.logger.WithField("race_id", race.RaceID)
	// 	_, err = s.universe.CreateRace(ctx, race)
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to insert race into db")
	// 		continue
	// 	}

	// 	err = s.cache.SetRace(ctx, race, cache.ExpiryHours(0))
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to cache race")
	// 	}

	// 	bar.Increment()
	// 	time.Sleep(time.Millisecond * 500)
	// }

	// bloodlines, _, err := s.esi.GetUniverseBloodlines(ctx, []*athena.Bloodline{})
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch bloodlines from ESI: %w", err)
	// }

	// name = "Bloodlines"
	// bar = p.AddBar(int64(len(bloodlines)),
	// 	mpb.PrependDecorators(
	// 		// simple name decorator
	// 		decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
	// 		// decor.DSyncWidth bit enables column width synchronization
	// 		decor.CountersNoUnit("%d / %d"),
	// 	),
	// 	mpb.AppendDecorators(decor.Percentage()),
	// )

	// for _, bloodline := range bloodlines {
	// 	entry := s.logger.WithField("bloodline_id", bloodline.BloodlineID)
	// 	_, err = s.universe.CreateBloodline(ctx, bloodline)
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to create bloodline in DB")
	// 		continue
	// 	}

	// 	err = s.cache.SetBloodline(ctx, bloodline, cache.ExpiryHours(0))
	// 	if err != nil {
	// 		s.logger.WithError(err).WithField("bloodline_id", bloodline.BloodlineID).Error("failed to cache bloodline")
	// 	}

	// 	bar.Increment()

	// 	time.Sleep(time.Millisecond * 50)
	// }

	// ancestries, _, err := s.esi.GetUniverseAncestries(ctx, []*athena.Ancestry{})
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch ancestries from ESI: %w", err)
	// }

	// name = "Ancestries"
	// bar = p.AddBar(int64(len(ancestries)),
	// 	mpb.PrependDecorators(
	// 		// simple name decorator
	// 		decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
	// 		// decor.DSyncWidth bit enables column width synchronization
	// 		decor.CountersNoUnit("%d / %d"),
	// 	),
	// 	mpb.AppendDecorators(decor.Percentage()),
	// )

	// for _, ancestry := range ancestries {
	// 	entry := s.logger.WithField("ancestry_id", ancestry.AncestryID)
	// 	_, err = s.universe.CreateAncestry(ctx, ancestry)
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to create ancestry in DB")
	// 		continue
	// 	}

	// 	err = s.cache.SetAncestry(ctx, ancestry, cache.ExpiryHours(0))
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to cache ancestry")
	// 	}

	// 	bar.Increment()
	// 	time.Sleep(time.Millisecond * 50)
	// }

	// factions, _, err := s.esi.GetUniverseFactions(ctx, []*athena.Faction{})
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch factions from ESI: %w", err)
	// }

	// name = "Factions"
	// bar = p.AddBar(int64(len(factions)),
	// 	mpb.PrependDecorators(
	// 		// simple name decorator
	// 		decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
	// 		// decor.DSyncWidth bit enables column width synchronization
	// 		decor.CountersNoUnit("%d / %d"),
	// 	),
	// 	mpb.AppendDecorators(decor.Percentage()),
	// )

	// for _, faction := range factions {
	// 	entry := s.logger.WithField("faction_id", faction.FactionID)
	// 	_, err = s.universe.CreateFaction(ctx, faction)
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to create faction in DB")
	// 		continue
	// 	}

	// 	err = s.cache.SetFaction(ctx, faction, cache.ExpiryHours(0))
	// 	if err != nil {
	// 		entry.WithError(err).Error("failed to cache faction")
	// 	}

	// 	bar.Increment()

	// 	time.Sleep(time.Millisecond * 50)
	// }

	categoryIDs, _, err := s.esi.GetUniverseCategories(ctx, []int{})
	if err != nil {
		return fmt.Errorf("failed to fetch category IDs from ESI: %w", err)
	}

	categoriesName := "Categories"
	categoriesBar := p.AddBar(int64(len(categoryIDs)),
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(categoriesName, decor.WC{W: len(categoriesName) + 1, C: decor.DidentRight}),
			// decor.DSyncWidth bit enables column width synchronization
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	for _, categoryID := range categoryIDs {
		categoryEntry := s.logger.WithField("category_id", categoryID)

		category, _, err := s.esi.GetUniverseCategoriesCategoryID(ctx, &athena.Category{CategoryID: categoryID})
		if err != nil {
			categoryEntry.WithError(err).Error("failed to fetch category from ESI")
			continue
		}

		_, err = s.universe.CreateCategory(ctx, category)
		if err != nil {
			categoryEntry.WithError(err).Error("failed to create category in DB")
			continue
		}

		groupsName := fmt.Sprintf("Groups of Category %d", categoryID)
		groupsBar := p.AddBar(int64(len(category.Groups)), mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(groupsName, decor.WC{W: len(groupsName) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, groupID := range category.Groups {
			wg.Add(1)
			go func(groupID int) {
				defer groupsBar.Increment()
				defer wg.Done()

				var wg2 = &sync.WaitGroup{}
				groupEntry := categoryEntry.WithField("group_id", groupID)

				group, _, err := s.esi.GetUniverseGroupsGroupID(ctx, &athena.Group{GroupID: groupID})
				if err != nil {
					groupEntry.WithError(err).Error("failed to fetch group from ESI")
					return
				}

				_, err = s.universe.CreateGroup(ctx, group)
				if err != nil {
					groupEntry.WithError(err).Error("failed to create group from DB")
					return
				}

				typesName := fmt.Sprintf("Types for Group %d", groupID)
				typesBar := p.AddBar(int64(len(group.Types)), mpb.BarRemoveOnComplete(),
					mpb.PrependDecorators(
						// simple name decorator
						decor.Name(typesName, decor.WC{W: len(typesName) + 1, C: decor.DidentRight}),
						// decor.DSyncWidth bit enables column width synchronization
						decor.CountersNoUnit("%d / %d"),
					),
					mpb.AppendDecorators(decor.Percentage()),
				)

				chunks := s.chunkSliceInts(group.Types, 50)

				for _, chunk := range chunks {
					wg2.Add(1)
					go func(typesBar *mpb.Bar, chunk []int, wg2 *sync.WaitGroup, groupEntry *logrus.Entry) {
						defer wg2.Done()
						for _, typeID := range chunk {
							typeEntry := groupEntry.WithField("type_id", typeID)
							item, _, err := s.esi.GetUniverseTypesTypeID(ctx, &athena.Type{TypeID: typeID})
							if err != nil {
								typeEntry.WithError(err).Error("failed to fetch type from ESI")
								continue
							}

							_, err = s.universe.CreateType(ctx, item)
							if err != nil {
								typeEntry.WithError(err).Error("failed to create type from DB")
								continue
							}

							err = s.cache.SetType(ctx, item, cache.ExpiryHours(0))
							if err != nil {
								typeEntry.WithError(err).Error("failed to cache type")
							}

							typesBar.Increment()

						}
					}(typesBar, chunk, wg2, groupEntry)
				}
				wg2.Wait()
				typesBar.SetTotal(int64(len(group.Types)), true)

				group.Types = nil

				err = s.cache.SetGroup(ctx, group, cache.ExpiryHours(0))
				if err != nil {
					groupEntry.WithError(err).Error("failed to cache category")
					return
				}

			}(groupID)

		}
		wg.Wait()
		groupsBar.SetTotal(int64(len(category.Groups)), true)

		category.Groups = nil

		err = s.cache.SetCategory(ctx, category, cache.ExpiryHours(0))
		if err != nil {
			categoryEntry.WithError(err).Error("failed to cache category")
		}

		categoriesBar.Increment()

	}

	categoriesBar.SetTotal(int64(len(categoryIDs)), true)

	// regionIDs, _, err := s.esi.GetUniverseRegions(ctx, []int{})
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch region id for ESI: %w", err)
	// }

	// regionsName := "Constellations"
	// regionsBar := p.AddBar(int64(len(regionIDs)),
	// 	mpb.PrependDecorators(
	// 		// simple name decorator
	// 		decor.Name(regionsName, decor.WC{W: len(regionsName) + 1, C: decor.DidentRight}),
	// 		// decor.DSyncWidth bit enables column width synchronization
	// 		decor.CountersNoUnit("%d / %d"),
	// 	),
	// 	mpb.AppendDecorators(decor.Percentage()),
	// )

	// for _, regionID := range regionIDs {

	// 	regionEntry := s.logger.WithField("region_id", regionID)

	// 	region, _, err := s.esi.GetUniverseRegionsRegionID(ctx, &athena.Region{RegionID: regionID})
	// 	if err != nil {
	// 		regionEntry.WithError(err).Error("failed to fetch region from ESI")
	// 		continue
	// 	}

	// 	_, err = s.universe.CreateRegion(ctx, region)
	// 	if err != nil {
	// 		regionEntry.WithError(err).Error("failed to create region in DB")
	// 		continue
	// 	}

	// 	constellationsName := "Constellations"
	// 	constellationsBar := p.AddBar(int64(len(region.ConstellationIDs)), mpb.BarRemoveOnComplete(),
	// 		mpb.PrependDecorators(
	// 			// simple name decorator
	// 			decor.Name(constellationsName, decor.WC{W: len(constellationsName) + 1, C: decor.DidentRight}),
	// 			// decor.DSyncWidth bit enables column width synchronization
	// 			decor.CountersNoUnit("%d / %d"),
	// 		),
	// 		mpb.AppendDecorators(decor.Percentage()),
	// 	)

	// 	for _, constellationID := range region.ConstellationIDs {
	// 		wg.Add(1)
	// 		go func(constellationID int) {
	// 			defer wg.Done()
	// 			defer constellationsBar.Increment()

	// 			constellationsEntry := regionEntry.WithField("constellation_id", constellationID)

	// 			constellation, _, err := s.esi.GetUniverseConstellationsConstellationID(ctx, &athena.Constellation{ConstellationID: constellationID})
	// 			if err != nil {
	// 				constellationsEntry.WithError(err).Error("failed to fetch constellation from ESI")
	// 				return
	// 			}

	// 			_, err = s.universe.CreateConstellation(ctx, constellation)
	// 			if err != nil {
	// 				constellationsEntry.WithError(err).Error("failed to create constellation in DB")
	// 				return
	// 			}

	// 			systemsName := fmt.Sprintf("Systems for Constellation %d", constellationID)
	// 			systemsBar := p.AddBar(int64(len(constellation.SystemIDs)), mpb.BarRemoveOnComplete(),
	// 				mpb.PrependDecorators(
	// 					// simple name decorator
	// 					decor.Name(systemsName, decor.WC{W: len(systemsName) + 1, C: decor.DidentRight}),
	// 					// decor.DSyncWidth bit enables column width synchronization
	// 					decor.CountersNoUnit("%d / %d"),
	// 				),
	// 				mpb.AppendDecorators(decor.Percentage()),
	// 			)

	// 			for _, systemID := range constellation.SystemIDs {
	// 				systemEntry := constellationsEntry.WithField("type_id", systemID)
	// 				system, _, err := s.esi.GetUniverseSolarSystemsSolarSystemID(ctx, &athena.SolarSystem{SystemID: systemID})
	// 				if err != nil {
	// 					systemEntry.WithError(err).Error("failed to fetch system from ESI")
	// 					continue
	// 				}

	// 				_, err = s.universe.CreateSolarSystem(ctx, system)
	// 				if err != nil {
	// 					systemEntry.WithError(err).Error("failed to create system from DB")
	// 					continue
	// 				}

	// 				err = s.cache.SetSolarSystem(ctx, system, cache.ExpiryHours(0))
	// 				if err != nil {
	// 					systemEntry.WithError(err).Error("failed to cache system")
	// 				}

	// 				systemsBar.Increment()

	// 			}

	// 			systemsBar.SetTotal(int64(len(constellation.SystemIDs)), true)

	// 			constellation.SystemIDs = nil

	// 			err = s.cache.SetConstellation(ctx, constellation, cache.ExpiryHours(0))
	// 			if err != nil {
	// 				constellationsEntry.WithError(err).Error("failed to cache constellation")
	// 				return
	// 			}
	// 		}(constellationID)
	// 	}
	// 	wg.Wait()
	// 	constellationsBar.SetTotal(int64(len(region.ConstellationIDs)), true)

	// 	region.ConstellationIDs = nil

	// 	err = s.cache.SetRegion(ctx, region, cache.ExpiryHours(0))
	// 	if err != nil {
	// 		regionEntry.WithError(err).Error("failed to cache region")
	// 	}

	// 	regionsBar.Increment()
	// }

	return nil
}

func (s *service) Ancestry(ctx context.Context, id int) (*athena.Ancestry, error) {

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

	err = s.cache.SetAncestry(ctx, ancestry, cache.ExpiryMinutes(0))

	return ancestry, err

}

func (s *service) Bloodline(ctx context.Context, id int) (*athena.Bloodline, error) {

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

	err = s.cache.SetBloodline(ctx, bloodline, cache.ExpiryMinutes(0))

	return bloodline, err

}

func (s *service) Category(ctx context.Context, id int) (*athena.Category, error) {

	category, err := s.cache.Category(ctx, id)
	if err != nil {
		return nil, err
	}

	if category != nil {
		return category, nil
	}

	category, err = s.universe.Category(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		category, _, err = s.esi.GetUniverseCategoriesCategoryID(ctx, &athena.Category{CategoryID: id})
		if err != nil {
			return nil, err
		}

		category, err = s.universe.CreateCategory(ctx, category)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetCategory(ctx, category, cache.ExpiryMinutes(0))

	return category, err

}

func (s *service) Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error) {
	return s.universe.Categories(ctx, operators...)
}

func (s *service) Constellation(ctx context.Context, id int) (*athena.Constellation, error) {

	constellation, err := s.cache.Constellation(ctx, id)
	if err != nil {
		return nil, err
	}

	if constellation != nil {
		return constellation, nil
	}

	constellation, err = s.universe.Constellation(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		constellation, _, err = s.esi.GetUniverseConstellationsConstellationID(ctx, &athena.Constellation{ConstellationID: id})
		if err != nil {
			return nil, err
		}

		constellation, err = s.universe.CreateConstellation(ctx, constellation)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetConstellation(ctx, constellation, cache.ExpiryMinutes(60))

	return constellation, err

}

func (s *service) Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error) {
	panic("universe.Constellations has not been implemented")
}

func (s *service) Faction(ctx context.Context, id int) (*athena.Faction, error) {
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

	err = s.cache.SetFaction(ctx, faction, cache.ExpiryMinutes(0))

	return faction, err
}

func (s *service) Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error) {
	panic("universe.Factions has not been implemented")
}

func (s *service) Group(ctx context.Context, id int) (*athena.Group, error) {

	group, err := s.cache.Group(ctx, id)
	if err != nil {
		return nil, err
	}

	if group != nil {
		return group, nil
	}

	group, err = s.universe.Group(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		group, _, err = s.esi.GetUniverseGroupsGroupID(ctx, &athena.Group{GroupID: id})
		if err != nil {
			return nil, err
		}

		group, err = s.universe.CreateGroup(ctx, group)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetGroup(ctx, group, cache.ExpiryMinutes(60))

	return group, err

}

func (s *service) Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error) {
	panic("universe.Groups has not been implemented")
}

func (s *service) Race(ctx context.Context, id int) (*athena.Race, error) {

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

	err = s.cache.SetRace(ctx, race, cache.ExpiryMinutes(0))

	return race, err

}

func (s *service) Region(ctx context.Context, id int) (*athena.Region, error) {

	region, err := s.cache.Region(ctx, id)
	if err != nil {
		return nil, err
	}

	if region != nil {
		return region, nil
	}

	region, err = s.universe.Region(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		region, _, err = s.esi.GetUniverseRegionsRegionID(ctx, &athena.Region{RegionID: id})
		if err != nil {
			return nil, err
		}

		region, err = s.universe.CreateRegion(ctx, region)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetRegion(ctx, region, cache.ExpiryMinutes(0))

	return region, err

}

func (s *service) Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error) {
	panic("universe.Regions has not been implemented")
}

func (s *service) SolarSystem(ctx context.Context, id int) (*athena.SolarSystem, error) {

	solarSystem, err := s.cache.SolarSystem(ctx, id)
	if err != nil {
		return nil, err
	}

	if solarSystem != nil {
		return solarSystem, nil
	}

	solarSystem, err = s.universe.SolarSystem(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		solarSystem, _, err = s.esi.GetUniverseSolarSystemsSolarSystemID(ctx, &athena.SolarSystem{SystemID: id})
		if err != nil {
			return nil, err
		}

		solarSystem, err = s.universe.CreateSolarSystem(ctx, solarSystem)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetSolarSystem(ctx, solarSystem, cache.ExpiryMinutes(30))

	return solarSystem, err

}

func (s *service) SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error) {
	panic("universe.SolarSystems has not been implemented")
}

func (s *service) Station(ctx context.Context, id int) (*athena.Station, error) {

	station, err := s.cache.Station(ctx, id)
	if err != nil {
		return nil, err
	}

	if station != nil {
		return station, nil
	}

	station, err = s.universe.Station(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		station, _, err = s.esi.GetUniverseStationsStationID(ctx, &athena.Station{StationID: id})
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

func (s *service) Structure(ctx context.Context, member *athena.Member, id int64) (*athena.Structure, error) {

	structure, err := s.cache.Structure(ctx, id)
	if err != nil {
		return nil, err
	}

	if structure != nil {
		return structure, nil
	}

	structure, err = s.universe.Structure(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		// TODO: Deliver a Concreate Error from ESI Package and insert th is into
		structure, _, err = s.esi.GetUniverseStructuresStructureID(ctx, member, &athena.Structure{StructureID: id})
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

func (s *service) Type(ctx context.Context, id int) (*athena.Type, error) {

	item, err := s.cache.Type(ctx, id)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return item, nil
	}

	item, err = s.universe.Type(ctx, id)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		item, _, err = s.esi.GetUniverseTypesTypeID(ctx, &athena.Type{TypeID: id})
		if err != nil {
			return nil, err
		}

		item, err = s.universe.CreateType(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetType(ctx, item, cache.ExpiryMinutes(30))

	return item, err

}

func (s *service) Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error) {
	panic("universe.Types has not been implemented")
}

func (s *service) chunkSliceInts(slc []int, size int) [][]int {
	var slcLen = len(slc)
	var divided = make([][]int, slcLen/size)

	for i := 0; i < slcLen; i += size {
		end := i + size

		if end > slcLen {
			end = slcLen
		}

		divided = append(divided, slc[i:end])
	}

	return divided
}
