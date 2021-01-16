package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/generated"
)

func (r *authAttemptResolver) Status(ctx context.Context, obj *athena.AuthAttempt) (string, error) {
	return obj.Status.String(), nil
}

// AuthAttempt returns generated.AuthAttemptResolver implementation.
func (r *Resolver) AuthAttempt() generated.AuthAttemptResolver { return &authAttemptResolver{r} }

type authAttemptResolver struct{ *Resolver }
