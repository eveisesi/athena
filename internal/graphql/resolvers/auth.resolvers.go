package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *authAttemptResolver) Status(ctx context.Context, obj *athena.AuthAttempt) (string, error) {
	return obj.Status.String(), nil
}

func (r *authAttemptResolver) URL(ctx context.Context, obj *athena.AuthAttempt, scopes []string) (string, error) {
	return r.auth.AuthorizationURI(ctx, obj.State, scopes), nil
}

// AuthAttempt returns service.AuthAttemptResolver implementation.
func (r *resolver) AuthAttempt() service.AuthAttemptResolver { return &authAttemptResolver{r} }

type authAttemptResolver struct{ *resolver }
