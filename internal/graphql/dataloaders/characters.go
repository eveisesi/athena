package dataloaders

import (
	"context"
	"sort"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/graphql/dataloaders/generated"
)

type characterLoaders struct {
	Character                   *generated.CharacterLoader
	CharacterCorporationHistory *generated.CharacterCorporationHistoryLoader
}

func newCharacterLoaders(ctx context.Context, c character.Service) *characterLoaders {
	return &characterLoaders{
		Character:                   characterLoader(ctx, c),
		CharacterCorporationHistory: characterCorporationHistoryLoader(ctx, c),
	}
}

func characterLoader(ctx context.Context, c character.Service) *generated.CharacterLoader {
	return generated.NewCharacterLoader(generated.CharacterLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([]*athena.Character, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([]*athena.Character, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := c.Characters(ctx, athena.NewInOperator("id", k))
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

func characterCorporationHistoryLoader(ctx context.Context, c character.Service) *generated.CharacterCorporationHistoryLoader {
	return generated.NewCharacterCorporationHistoryLoader(generated.CharacterCorporationHistoryLoaderConfig{
		Wait:     defaultWait,
		MaxBatch: defaultMaxBatch,
		Fetch: func(keys []uint) ([][]*athena.CharacterCorporationHistory, []error) {
			var errors = make([]error, 0, len(keys))
			var results = make([][]*athena.CharacterCorporationHistory, len(keys))

			k := append(make([]uint, 0, len(keys)), keys...)
			sort.SliceStable(k, func(i, j int) bool {
				return k[i] < k[j]
			})

			rows, err := c.CharacterCorporationHistory(ctx, athena.NewInOperator("character_id", k))
			if err != nil {
				errors = append(errors, err)
				return nil, errors
			}

			resultsByPrimaryKey := make(map[uint][]*athena.CharacterCorporationHistory)
			for _, row := range rows {
				resultsByPrimaryKey[row.CharacterID] = append(resultsByPrimaryKey[row.CharacterID], row)
			}

			for i, v := range keys {
				results[i] = resultsByPrimaryKey[v]
			}

			return results, nil

		},
	})
}
