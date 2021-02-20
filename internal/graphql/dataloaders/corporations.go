package dataloaders

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
)

func corporationLoader(ctx context.Context, c corporation.Service) *generated.CorporationLoader {
	return generated.NewCorporationLoader(generated.CorporationLoaderConfig{
		Fetch: func(keys []uint) ([]*athena.Corporation, []error) {

			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Corporation, len(keys))

			rows, err := c.Corporations(ctx, athena.NewOperators(athena.NewInOperator("id", keys)))
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
