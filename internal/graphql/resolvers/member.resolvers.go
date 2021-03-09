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

func (r *memberResolver) Scopes(ctx context.Context, obj *athena.Member) ([]string, error) {
	s := make([]string, 0, len(obj.Scopes))
	for _, i := range obj.Scopes {
		s = append(s, i.Scope.String())
	}

	return s, nil
}

func (r *memberResolver) Main(ctx context.Context, obj *athena.Member) (*athena.Character, error) {
	if !obj.MainID.Valid {
		return nil, nil
	}

	return dataloaders.CtxLoaders(ctx).Character.Load(obj.MainID.Uint)
}

func (r *memberResolver) Character(ctx context.Context, obj *athena.Member) (*athena.Character, error) {
	return dataloaders.CtxLoaders(ctx).Character.Load(obj.ID)
}

func (r *queryResolver) Member(ctx context.Context) (*athena.Member, error) {
	// r.member.MemberFromContext(ctx)

	return nil, fmt.Errorf("Failed to fetch member from context")
}

// Member returns service.MemberResolver implementation.
func (r *resolver) Member() service.MemberResolver { return &memberResolver{r} }

type memberResolver struct{ *resolver }
