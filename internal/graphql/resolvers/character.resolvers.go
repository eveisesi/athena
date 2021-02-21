package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/dataloaders"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *characterResolver) Race(ctx context.Context, obj *athena.Character) (*athena.Race, error) {
	return dataloaders.CtxLoaders(ctx).Race.Load(obj.RaceID)
}

func (r *characterResolver) Bloodline(ctx context.Context, obj *athena.Character) (*athena.Bloodline, error) {
	return dataloaders.CtxLoaders(ctx).Bloodline.Load(obj.BloodlineID)
}

func (r *characterResolver) Ancestry(ctx context.Context, obj *athena.Character) (*athena.Ancestry, error) {
	if !obj.AncestryID.Valid {
		return nil, nil
	}
	return dataloaders.CtxLoaders(ctx).Ancestry.Load(obj.AncestryID.Uint)
}

func (r *characterResolver) Corporation(ctx context.Context, obj *athena.Character) (*athena.Corporation, error) {
	return dataloaders.CtxLoaders(ctx).Corporation.Load(obj.CorporationID)
}

func (r *characterResolver) Alliance(ctx context.Context, obj *athena.Character) (*athena.Alliance, error) {
	if !obj.AllianceID.Valid {
		return nil, nil
	}
	return dataloaders.CtxLoaders(ctx).Alliance.Load(obj.AllianceID.Uint)
}

// Character returns service.CharacterResolver implementation.
func (r *resolver) Character() service.CharacterResolver { return &characterResolver{r} }

type characterResolver struct{ *resolver }
