package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *memberAssetResolver) LocationFlag(ctx context.Context, obj *athena.MemberAsset) (string, error) {
	return obj.LocationFlag.String(), nil
}

func (r *memberAssetResolver) LocationType(ctx context.Context, obj *athena.MemberAsset) (string, error) {
	return obj.LocationType.String(), nil
}

func (r *queryResolver) MemberAssets(ctx context.Context, memberID uint, page uint) ([]*athena.MemberAsset, error) {
	return r.asset.MemberAssets(ctx, memberID, page)
}

// MemberAsset returns service.MemberAssetResolver implementation.
func (r *resolver) MemberAsset() service.MemberAssetResolver { return &memberAssetResolver{r} }

type memberAssetResolver struct{ *resolver }
