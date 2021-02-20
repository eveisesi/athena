package resolvers

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql"
)

func (r *memberResolver) OwnerHash(ctx context.Context, obj *athena.Member) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberResolver) Scopes(ctx context.Context, obj *athena.Member) ([]*athena.MemberScope, error) {
	s := make([]*athena.MemberScope, 0, len(obj.Scopes))
	for _, i := range obj.Scopes {
		s = append(s, &i)
	}

	return s, nil
}

func (r *memberScopeResolver) Scope(ctx context.Context, obj *athena.MemberScope) (string, error) {
	return obj.Scope.String(), nil
}

func (r *queryResolver) Member(ctx context.Context) (*athena.Member, error) {
	member := r.member.MemberFromContext(ctx)

	return member, nil
}

// Member returns graphql.MemberResolver implementation.
func (r *resolver) Member() graphql.MemberResolver { return &memberResolver{r} }

// MemberScope returns generated.MemberScopeResolver implementation.
func (r *resolver) MemberScope() graphql.MemberScopeResolver { return &memberScopeResolver{r} }

type memberResolver struct{ *resolver }
type memberScopeResolver struct{ *resolver }
