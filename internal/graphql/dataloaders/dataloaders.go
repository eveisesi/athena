package dataloaders

import (
	"context"
	"net/http"
	"time"

	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/universe"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{name: "dataloaders"}

const (
	defaultWait     = 100 * time.Millisecond
	defaultMaxBatch = 500
)

type Loaders struct {
	*allianceLoaders
	// *characterLoaders
	// *corporationLoaders
	// *universeLoaders
}

func Middleware(a alliance.Service, ch character.Service, corp corporation.Service, uni universe.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		ctx = context.WithValue(ctx, ctxKey, Loaders{
			allianceLoaders: newAllianceLoaders(ctx, a),
			// characterLoaders: newCharacterLoaders(ctx, ch),
			// newCorporationLoaders(ctx, corp),
			// newUniverseLoaders(ctx, uni),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
