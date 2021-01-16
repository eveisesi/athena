package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/generated"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (r *queryResolver) Auth(ctx context.Context) (*athena.AuthAttempt, error) {
	attempt, err := r.auth.InitializeAttempt(ctx)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	return attempt, nil
}

func (r *queryResolver) AuthorizationURI(ctx context.Context, state string) (string, error) {
	return r.auth.AuthorizationURI(ctx, state), nil
}

func (r *subscriptionResolver) AuthStatus(ctx context.Context, state string) (<-chan *athena.AuthAttempt, error) {
	pipe := make(chan *athena.AuthAttempt)

	go func(ctx context.Context, pipe chan *athena.AuthAttempt, state string) {
		ctx = newrelic.NewContext(ctx, newrelic.FromContext(ctx).NewGoroutine())

		ticker := time.NewTicker(time.Second * 2)

		for {

			select {
			case <-ctx.Done():
				close(pipe)
				return
			case <-ticker.C:
				attempt, err := r.auth.AuthAttempt(ctx, state)
				if err != nil {
					newrelic.FromContext(ctx).NoticeError(err)
					ticker.Stop()
					close(pipe)
					return
				}

				pipe <- attempt

				if attempt.Status != athena.PendingAuthStatus {
					ticker.Stop()
					close(pipe)
					return
				}
				break
			}

		}

	}(ctx, pipe, state)

	return pipe, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }