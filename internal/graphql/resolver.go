package graphql

import (
	"github.com/eveisesi/athena/internal/auth"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger *logrus.Logger

	auth auth.Service
}

func NewResolvers(logger *logrus.Logger, auth auth.Service) *Resolver {
	return &Resolver{
		logger,
		auth,
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
