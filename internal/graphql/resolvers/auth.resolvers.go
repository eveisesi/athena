package resolvers

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql"
)

func (r *authAttemptResolver) Status(ctx context.Context, obj *athena.AuthAttempt) (string, error) {
	return obj.Status.String(), nil
}

func (r *authAttemptResolver) URL(ctx context.Context, obj *athena.AuthAttempt, scopes []string) (string, error) {
	return r.auth.AuthorizationURI(ctx, obj.State, scopes), nil
}

// AuthAttempt returns generated.AuthAttemptResolver implementation.
func (r *resolver) AuthAttempt() graphql.AuthAttemptResolver { return &authAttemptResolver{r} }

type authAttemptResolver struct{ *resolver }
