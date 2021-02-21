package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/dataloaders"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *memberHomeLocationResolver) Info(ctx context.Context, obj *athena.MemberHomeLocation) (service.CloneLocationInfo, error) {
	switch obj.LocationType {
	case "station":
		return dataloaders.CtxLoaders(ctx).Station.Load(uint(obj.LocationID))
	case "structure":
		return dataloaders.CtxLoaders(ctx).Structure.Load(obj.LocationID)
	default:
		return nil, fmt.Errorf("%v is not a resolvable location type", obj.LocationType)
	}
}

func (r *memberImplantResolver) Type(ctx context.Context, obj *athena.MemberImplant) (*athena.Type, error) {
	return dataloaders.CtxLoaders(ctx).Item.Load(obj.ImplantID)
}

func (r *memberJumpCloneResolver) Implants(ctx context.Context, obj *athena.MemberJumpClone) ([]*athena.Type, error) {
	items := make([]*athena.Type, 0, len(obj.Implants))
	for _, t := range obj.Implants {
		item, e := dataloaders.CtxLoaders(ctx).Item.Load(uint(t))
		if e != nil {
			return nil, e
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *memberJumpCloneResolver) Info(ctx context.Context, obj *athena.MemberJumpClone) (service.CloneLocationInfo, error) {
	switch obj.LocationType {
	case "station":
		return dataloaders.CtxLoaders(ctx).Station.Load(uint(obj.LocationID))
	case "structure":
		return dataloaders.CtxLoaders(ctx).Structure.Load(obj.LocationID)
	default:
		return nil, fmt.Errorf("%v is not a resolvable location type", obj.LocationType)
	}
}

func (r *queryResolver) MemberClones(ctx context.Context, memberID uint) (*athena.MemberClones, error) {
	return r.clone.MemberClones(ctx, memberID)
}

func (r *queryResolver) MemberImplants(ctx context.Context, memberID uint) ([]*athena.MemberImplant, error) {
	return r.clone.MemberImplants(ctx, memberID)
}

// MemberHomeLocation returns service.MemberHomeLocationResolver implementation.
func (r *resolver) MemberHomeLocation() service.MemberHomeLocationResolver {
	return &memberHomeLocationResolver{r}
}

// MemberImplant returns service.MemberImplantResolver implementation.
func (r *resolver) MemberImplant() service.MemberImplantResolver { return &memberImplantResolver{r} }

// MemberJumpClone returns service.MemberJumpCloneResolver implementation.
func (r *resolver) MemberJumpClone() service.MemberJumpCloneResolver {
	return &memberJumpCloneResolver{r}
}

type memberHomeLocationResolver struct{ *resolver }
type memberImplantResolver struct{ *resolver }
type memberJumpCloneResolver struct{ *resolver }
