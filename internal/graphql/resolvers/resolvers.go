package resolvers

import (
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/asset"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/contact"
	"github.com/eveisesi/athena/internal/contract"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/graphql/service"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type resolver struct {
	logger *logrus.Logger

	auth        auth.Service
	member      member.Service
	character   character.Service
	corporation corporation.Service
	alliance    alliance.Service
	universe    universe.Service
	location    location.Service
	clone       clone.Service
	contact     contact.Service
	contract    contract.Service
	asset       asset.Service
}

func New(
	logger *logrus.Logger,
	auth auth.Service,
	member member.Service,
	character character.Service,
	corporation corporation.Service,
	alliance alliance.Service,
	universe universe.Service,
	location location.Service,
	clone clone.Service,
	contact contact.Service,
	contract contract.Service,
	asset asset.Service,
) service.ResolverRoot {
	return &resolver{
		logger:      logger,
		auth:        auth,
		member:      member,
		character:   character,
		corporation: corporation,
		alliance:    alliance,
		universe:    universe,
		clone:       clone,
		location:    location,
		contact:     contact,
		contract:    contract,
		asset:       asset,
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
