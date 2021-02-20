package dataloaders

import (
	"context"

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

			k := make([]interface{}, 0, len(keys))
			for _, key := range keys {
				k = append(k, key)
			}

			rows, err := a.Alliances(ctx, athena.NewInOperator("id", k...))
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
