package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/service"
)

func (r *structureResolver) ID(ctx context.Context, obj *athena.Structure) (uint, error) {
	panic(fmt.Errorf("not implemented"))
}

// Structure returns service.StructureResolver implementation.
func (r *resolver) Structure() service.StructureResolver { return &structureResolver{r} }

type structureResolver struct{ *resolver }
