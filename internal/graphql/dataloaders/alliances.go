package dataloaders

import (
	"context"
	"sort"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
)

type allianceLoaders struct {
	Alliance *generated.AllianceLoader
}

func newAllianceLoaders(ctx context.Context, a alliance.Service) *allianceLoaders {
	return &allianceLoaders{
		Alliance: allianceLoader(ctx, a),
	}
}

func allianceLoader(ctx context.Context, a alliance.Service) *generated.AllianceLoader {
	return generated.NewAllianceLoader(generated.AllianceLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Alliance, []error) {

			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Alliance, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := a.Alliances(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Alliance)
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
