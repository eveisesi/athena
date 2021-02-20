package resolvers

import (
	"github.com/eveisesi/athena/internal/graphql"
)

type queryResolver struct{ *resolver }

// Query returns generated.QueryResolver implementation.
func (r *resolver) Query() graphql.QueryResolver { return &queryResolver{r} }

type subscriptionResolver struct{ *resolver }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *resolver) Subscription() graphql.SubscriptionResolver { return &subscriptionResolver{r} }
