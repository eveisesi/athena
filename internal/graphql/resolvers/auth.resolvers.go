package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/graphql/service"
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

func (r *authAttemptResolver) Status(ctx context.Context, obj *athena.AuthAttempt) (string, error) {
	return obj.Status.String(), nil
}

func (r *authAttemptResolver) URL(ctx context.Context, obj *athena.AuthAttempt, scopes []string) (string, error) {
	return r.auth.AuthorizationURI(ctx, obj.State, scopes), nil
}

// AuthAttempt returns service.AuthAttemptResolver implementation.
func (r *resolver) AuthAttempt() service.AuthAttemptResolver { return &authAttemptResolver{r} }

type authAttemptResolver struct{ *resolver }
