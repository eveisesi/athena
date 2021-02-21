package dataloaders

import (
	"context"
	"sort"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
	"github.com/eveisesi/athena/internal/universe"
)

type universeLoaders struct {
	Ancestry    *generated.AncestryLoader
	Bloodline   *generated.BloodlineLoader
	Faction     *generated.FactionLoader
	Race        *generated.RaceLoader
	Category    *generated.CategoryLoader
	Group       *generated.GroupLoader
	Item        *generated.TypeLoader
	SolarSystem *generated.SolarSystemLoader
	Station     *generated.StationLoader
	Structure   *generated.StructureLoader
}

func newUniverseLoader(ctx context.Context, u universe.Service) *universeLoaders {
	return &universeLoaders{
		Ancestry:    ancestryLoader(ctx, u),
		Bloodline:   bloodlineLoader(ctx, u),
		Faction:     factionLoader(ctx, u),
		Race:        raceLoader(ctx, u),
		Category:    categoryLoader(ctx, u),
		Group:       groupLoader(ctx, u),
		Item:        typeLoader(ctx, u),
		SolarSystem: solarSystemLoader(ctx, u),
		Station:     stationLoader(ctx, u),
		Structure:   structureLoader(ctx, u),
	}
}

func ancestryLoader(ctx context.Context, u universe.Service) *generated.AncestryLoader {
	return generated.NewAncestryLoader(generated.AncestryLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Ancestry, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Ancestry, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Ancestries(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Ancestry)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func bloodlineLoader(ctx context.Context, u universe.Service) *generated.BloodlineLoader {
	return generated.NewBloodlineLoader(generated.BloodlineLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Bloodline, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Bloodline, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Bloodlines(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Bloodline)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func factionLoader(ctx context.Context, u universe.Service) *generated.FactionLoader {
	return generated.NewFactionLoader(generated.FactionLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Faction, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Faction, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Factions(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Faction)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func raceLoader(ctx context.Context, u universe.Service) *generated.RaceLoader {
	return generated.NewRaceLoader(generated.RaceLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Race, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Race, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Races(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Race)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func categoryLoader(ctx context.Context, u universe.Service) *generated.CategoryLoader {
	return generated.NewCategoryLoader(generated.CategoryLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Category, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Category, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Categories(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Category)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func groupLoader(ctx context.Context, u universe.Service) *generated.GroupLoader {
	return generated.NewGroupLoader(generated.GroupLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Group, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Group, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Groups(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Group)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func typeLoader(ctx context.Context, u universe.Service) *generated.TypeLoader {
	return generated.NewTypeLoader(generated.TypeLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Type, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Type, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Types(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Type)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func structureLoader(ctx context.Context, u universe.Service) *generated.StructureLoader {
	return generated.NewStructureLoader(generated.StructureLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint64) ([]*athena.Structure, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Structure, len(keys))

			k := append(make([]uint64, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Structures(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint64]*athena.Structure)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func stationLoader(ctx context.Context, u universe.Service) *generated.StationLoader {
	return generated.NewStationLoader(generated.StationLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Station, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Station, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.Stations(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Station)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}

func solarSystemLoader(ctx context.Context, u universe.Service) *generated.SolarSystemLoader {
	return generated.NewSolarSystemLoader(generated.SolarSystemLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.SolarSystem, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.SolarSystem, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := u.SolarSystems(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.SolarSystem)
			for _, row := range rows {
				resultsByPrimaryKey[row.ID] = row
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}
