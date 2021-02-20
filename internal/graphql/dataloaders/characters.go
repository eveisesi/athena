package dataloaders

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
)

type characterLoaders struct {
	character *generated.CharacterLoader
	history   *generated.CharacterCorporationHistoryLoader
}

// func newCharacterLoaders(ctx context.Context, c character.Service) *characterLoaders {
// 	return &characterLoaders{
// 		character: characterLoader(ctx, c),
// 		history:
// 	}
// }

func characterLoader(ctx context.Context, c character.Service) *generated.CharacterLoader {
	return generated.NewCharacterLoader(generated.CharacterLoaderConfig{
		Fetch: func(keys []uint) ([]*athena.Character, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Character, len(keys))

			rows, err := c.Characters(ctx, athena.NewOperators(athena.NewInOperator("id", keys)))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint]*athena.Character)
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

// func characterCorporationHistoryLoader(ctx context.Context, c character.Service) *generated.CharacterCorporationHistoryLoader {
// 	return generated.NewCharacterCorporationHistoryLoader(generated.CharacterCorporationHistoryLoaderConfig{
// 		Fetch: func(keys []uint) ([][]*athena.CharacterCorporationHistory, []error) {
// 			var errors = make([]error, 0, len(keys))
// 			var results = make([]*athena.Character, len(keys))

// 			rows, err := c.CharacterCorporationHistory(ctx, athena.NewOperators(athena.NewInOperator("id", keys)))
// 			if err != nil {
// 				errors = append(errors, err)
// 				return nil, errors
// 			}

// 		},
// 	})
// }
