package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/eveisesi/athena"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/auth"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/contact"
	"github.com/eveisesi/athena/internal/contract"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/graphql/resolvers"
	graphql "github.com/eveisesi/athena/internal/graphql/service"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type server struct {
	port     uint
	env      athena.Environment
	logger   *logrus.Logger
	newrelic *newrelic.Application

	auth        auth.Service
	cache       cache.Service
	member      member.Service
	character   character.Service
	corporation corporation.Service
	alliance    alliance.Service
	universe    universe.Service
	location    location.Service
	clone       clone.Service
	contact     contact.Service
	contract    contract.Service

	server *http.Server
}

func NewServer(
	port uint,
	env athena.Environment,
	logger *logrus.Logger,
	cache cache.Service,
	newrelic *newrelic.Application,
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
) *server {

	s := &server{
		port:        port,
		env:         env,
		logger:      logger,
		cache:       cache,
		newrelic:    newrelic,
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
	}

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 5,
		ErrorLog:          log.New(logger.Writer(), "", 0),
		Handler:           s.buildRouter(),
	}

	return s

}

func (s *server) Run() error {
	s.logger.WithField("port", s.port).Info("starting http server")
	return s.server.ListenAndServe()
}

func (s *server) buildRouter() *chi.Mux {

	r := chi.NewRouter()

	r.Use(
		s.cors,
	// s.monitoring,
	)

	r.Get("/auth/callback", s.handleGetAuthCallback)
	// r.Get("/auth/login", s.handleGetAuthLogin)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.SetHeader("Content-Type", "application/json"),
			s.member.Middleware,
		)

		r.Group(func(r chi.Router) {
			r.Use(s.dataloaders)

			// directives := graphql.NewDirectives()
			es := graphql.NewExecutableSchema(graphql.Config{
				Resolvers: resolvers.New(
					s.logger, s.auth, s.member,
					s.character, s.corporation, s.alliance,
					s.universe, s.location, s.clone,
					s.contact, s.contract,
				),
				// Directives: generated.DirectiveRoot{HasGrant: directives.HasGrant},
			})
			queryHandler := handler.New(es)

			queryHandler.AddTransport(transport.Websocket{
				KeepAlivePingInterval: 2 * time.Second,
			})
			queryHandler.AddTransport(transport.POST{})

			queryHandler.SetQueryCache(lru.New(1000))

			queryHandler.Use(extension.AutomaticPersistedQuery{
				Cache: lru.New(100),
			})

			if s.env != athena.ProductionEnvironment {
				queryHandler.Use(extension.Introspection{})
				r.Handle("/", playground.Handler("GraphQL Playground", "/query"))
			}

			r.Handle("/query", queryHandler)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r

}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
	return s.server.Shutdown(ctx)
}

func (s *server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error) {

	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
	if err != nil {
		s.writeResponse(ctx, w, code, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	s.writeResponse(ctx, w, code, nil)

}
