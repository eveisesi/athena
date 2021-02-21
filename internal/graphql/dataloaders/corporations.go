package dataloaders

import (
	"context"
	"sort"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
)

type corporationLoaders struct {
	Corporation                *generated.CorporationLoader
	CorporationAllianceHistory *generated.CorporationAllianceHistoryLoader
}

func newCorporationLoaders(ctx context.Context, c corporation.Service) *corporationLoaders {
	return &corporationLoaders{
		Corporation:                corporationLoader(ctx, c),
		CorporationAllianceHistory: corporationAllianceHistoryLoader(ctx, c),
	}
}

func corporationLoader(ctx context.Context, c corporation.Service) *generated.CorporationLoader {
	return generated.NewCorporationLoader(generated.CorporationLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Corporation, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Corporation, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := c.Corporations(ctx, athena.NewInOperator("id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Corporation)
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

func corporationAllianceHistoryLoader(ctx context.Context, c corporation.Service) *generated.CorporationAllianceHistoryLoader {
	return generated.NewCorporationAllianceHistoryLoader(generated.CorporationAllianceHistoryLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([][]*athena.CorporationAllianceHistory, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([][]*athena.CorporationAllianceHistory, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := c.CorporationAllianceHistory(ctx, athena.NewInOperator("corporation_id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint][]*athena.CorporationAllianceHistory)
			for _, row := range rows {
				resultsByPrimaryKey[row.CorporationID] = append(resultsByPrimaryKey[row.CorporationID], row)
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil
		},
	})
}
