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

func (r *memberContractResolver) Availability(ctx context.Context, obj *athena.MemberContract) (string, error) {
	return obj.Availability.String(), nil
}

func (r *memberContractResolver) Status(ctx context.Context, obj *athena.MemberContract) (string, error) {
	return obj.Status.String(), nil
}

func (r *memberContractResolver) Type(ctx context.Context, obj *athena.MemberContract) (string, error) {
	return obj.Type.String(), nil
}

func (r *memberContractResolver) Items(ctx context.Context, obj *athena.MemberContract) ([]*athena.MemberContractItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberContractResolver) Bids(ctx context.Context, obj *athena.MemberContract) ([]*athena.MemberContractBid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberContractBidResolver) Bidder(ctx context.Context, obj *athena.MemberContractBid) (*athena.Character, error) {
	return dataloaders.CtxLoaders(ctx).Character.Load(obj.BidderID)
}

func (r *queryResolver) MemberContracts(ctx context.Context, memberID uint, page uint) ([]*athena.MemberContract, error) {
	return r.contract.MemberContracts(ctx, memberID, page)
}

// MemberContract returns service.MemberContractResolver implementation.
func (r *resolver) MemberContract() service.MemberContractResolver { return &memberContractResolver{r} }

// MemberContractBid returns service.MemberContractBidResolver implementation.
func (r *resolver) MemberContractBid() service.MemberContractBidResolver {
	return &memberContractBidResolver{r}
}

type memberContractResolver struct{ *resolver }
type memberContractBidResolver struct{ *resolver }
