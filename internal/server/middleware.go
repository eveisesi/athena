package server

import (
	"context"
	"net/http"

	"github.com/eveisesi/athena/internal/graphql/dataloaders"
)

// Cors middleware to allow frontend consumption
func (s *server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "600")

		// Just return for OPTIONS
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) dataloaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		ctx = context.WithValue(
			ctx,
			dataloaders.CtxKey,
			dataloaders.New(
				ctx,
				s.alliance,
				s.character,
				s.corporation,
				s.universe,
			),
		)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
