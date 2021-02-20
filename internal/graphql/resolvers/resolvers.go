package resolvers

import (
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/graphql"
	"github.com/eveisesi/athena/internal/member"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type resolver struct {
	logger *logrus.Logger

	auth   auth.Service
	member member.Service
}

func New(logger *logrus.Logger, auth auth.Service, member member.Service) graphql.ResolverRoot {
	return &resolver{
		logger,
		auth,
		member,
	}
}

// func NewDirectives() directives {
// 	return directives{
// 		HasGrant: hasGrant,
// 	}
// }

// type directives struct {
// 	HasGrant func(ctx context.Context, obj interface{}, next graphql.Resolver, scope string) (interface{}, error)
// }

// // TODO: Clean this up. It is a POC I typed up at ll:30 at night :-P
// func hasGrant(ctx context.Context, obj interface{}, next graphql.Resolver, scope string) (interface{}, error) {
// 	fieldCtx := graphql.GetFieldContext(ctx)
// 	if fieldCtx == nil {
// 		return nil, fmt.Errorf("failed to determine grant validity")
// 	}

// 	args := fieldCtx.Args
// 	if _, ok := args["state"]; !ok {
// 		return nil, fmt.Errorf("failed to determine grant validity due to missing query arguement")
// 	}

// 	if _, ok := args["state"].(string); !ok {
// 		return nil, fmt.Errorf("failed to determine grant validity. Expected int, got %T", args["state"])
// 	}

// 	state := args["state"].(string)

// 	spew.Dump(state, scope)

// 	return nil, nil

// }
