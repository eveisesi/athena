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

func (r *memberContactResolver) LabelIDs(ctx context.Context, obj *athena.MemberContact) ([]uint64, error) {
	return []uint64(obj.LabelIDs), nil
}

func (r *memberContactResolver) Info(ctx context.Context, obj *athena.MemberContact) (service.ContactInfo, error) {
	switch obj.ContactType {
	case "character":
		return dataloaders.CtxLoaders(ctx).Character.Load(obj.ContactID)
	case "corporation":
		return dataloaders.CtxLoaders(ctx).Corporation.Load(obj.ContactID)
	case "alliance":
		return dataloaders.CtxLoaders(ctx).Alliance.Load(obj.ContactID)
	case "faction":
		return dataloaders.CtxLoaders(ctx).Faction.Load(obj.ContactID)
	default:
		return nil, fmt.Errorf("%v is not a resolvable contact type", obj.ContactType)
	}
}

func (r *queryResolver) MemberContacts(ctx context.Context, memberID uint, page uint) ([]*athena.MemberContact, error) {
	return r.contact.MemberContacts(ctx, memberID, page)
}

// MemberContact returns service.MemberContactResolver implementation.
func (r *resolver) MemberContact() service.MemberContactResolver { return &memberContactResolver{r} }

type memberContactResolver struct{ *resolver }
