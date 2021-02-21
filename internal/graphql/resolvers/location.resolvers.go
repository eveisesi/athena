package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/dataloaders"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *memberLocationResolver) System(ctx context.Context, obj *athena.MemberLocation) (*athena.SolarSystem, error) {
	return dataloaders.CtxLoaders(ctx).SolarSystem.Load(obj.SolarSystemID)
}

func (r *memberLocationResolver) Station(ctx context.Context, obj *athena.MemberLocation) (*athena.Station, error) {
	if !obj.StationID.Valid {
		return nil, nil
	}
	return dataloaders.CtxLoaders(ctx).Station.Load(obj.StationID.Uint)
}

func (r *memberLocationResolver) Structure(ctx context.Context, obj *athena.MemberLocation) (*athena.Structure, error) {
	if !obj.StructureID.Valid {
		return nil, nil
	}
	return dataloaders.CtxLoaders(ctx).Structure.Load(obj.StructureID.Uint64)
}

func (r *memberShipResolver) Ship(ctx context.Context, obj *athena.MemberShip) (*athena.Type, error) {
	return dataloaders.CtxLoaders(ctx).Item.Load(obj.ShipTypeID)
}

func (r *queryResolver) MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error) {
	return r.location.MemberLocation(ctx, memberID)
}

func (r *queryResolver) MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error) {
	return r.location.MemberOnline(ctx, memberID)
}

func (r *queryResolver) MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error) {
	return r.location.MemberShip(ctx, memberID)
}

// MemberLocation returns service.MemberLocationResolver implementation.
func (r *resolver) MemberLocation() service.MemberLocationResolver { return &memberLocationResolver{r} }

// MemberShip returns service.MemberShipResolver implementation.
func (r *resolver) MemberShip() service.MemberShipResolver { return &memberShipResolver{r} }

type memberLocationResolver struct{ *resolver }
type memberShipResolver struct{ *resolver }
