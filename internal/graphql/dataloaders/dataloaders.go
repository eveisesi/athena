package dataloaders

import (
	"context"
	"time"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/universe"
)

type ctxKeyType struct{ name string }

var CtxKey = ctxKeyType{name: "dataloaders"}

const (
	defaultWait     = 100 * time.Millisecond
	defaultMaxBatch = 500
)

type Loaders struct {
	*allianceLoaders
	*characterLoaders
	*corporationLoaders
	*universeLoaders
}

func New(ctx context.Context, a alliance.Service, ch character.Service, corp corporation.Service, u universe.Service) *Loaders {

	return &Loaders{
		allianceLoaders:    newAllianceLoaders(ctx, a),
		characterLoaders:   newCharacterLoaders(ctx, ch),
		corporationLoaders: newCorporationLoaders(ctx, corp),
		universeLoaders:    newUniverseLoader(ctx, u),
	}

}

// CtxLoaders will extract available Loaders from the specified context
func CtxLoaders(ctx context.Context) *Loaders {
	return ctx.Value(CtxKey).(*Loaders)
}
