package dataloaders

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
	"github.com/eveisesi/athena/internal/universe"
)

type universeLoaders struct {
	Ancestry  *generated.AncestryLoader
	Bloodline *generated.BloodlineLoader
	Race      *generated.RaceLoader
	Category  *generated.CategoryLoader
	Group     *generated.GroupLoader
	Item      *generated.TypeLoader
}

func newUniverseLoader(ctx context.Context, u universe.Service) *universeLoaders {
	return &universeLoaders{
		Ancestry:  ancestryLoader(ctx, u),
		Bloodline: bloodlineLoader(ctx, u),
		Race:      raceLoader(ctx, u),
		Category:  categoryLoader(ctx, u),
		Group:     groupLoader(ctx, u),
		Item:      typeLoader(ctx, u),
	}
}

func ancestryLoader(ctx context.Context, u universe.Service) *generated.AncestryLoader {
	return generated.NewAncestryLoader(generated.AncestryLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Ancestry, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Ancestry, len(keys))

			rows, err := u.Ancestries(ctx, athena.NewInOperator("id", keys))
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

			rows, err := u.Bloodlines(ctx, athena.NewInOperator("id", keys))
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

func raceLoader(ctx context.Context, u universe.Service) *generated.RaceLoader {
	return generated.NewRaceLoader(generated.RaceLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Race, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Race, len(keys))

			rows, err := u.Races(ctx, athena.NewInOperator("id", keys))
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

			rows, err := u.Categories(ctx, athena.NewInOperator("id", keys))
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

			rows, err := u.Groups(ctx, athena.NewInOperator("id", keys))
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

			rows, err := u.Types(ctx, athena.NewInOperator("id", keys))
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
