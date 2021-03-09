package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/eveisesi/athena/internal/graphql/service"
)

// Query returns service.QueryResolver implementation.
func (r *resolver) Query() service.QueryResolver { return &queryResolver{r} }

// Subscription returns service.SubscriptionResolver implementation.
func (r *resolver) Subscription() service.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *resolver }
type subscriptionResolver struct{ *resolver }
