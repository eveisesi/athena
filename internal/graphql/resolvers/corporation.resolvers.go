package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *corporationResolver) Shares(ctx context.Context, obj *athena.Corporation) (uint, error) {
	panic(fmt.Errorf("not implemented"))
}

// Corporation returns service.CorporationResolver implementation.
func (r *resolver) Corporation() service.CorporationResolver { return &corporationResolver{r} }

type corporationResolver struct{ *resolver }
