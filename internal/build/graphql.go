package build

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cockroachdb/errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	graphqlv1 "github.com/powerdigital/go-micro/internal/transport/graphql/v1"
)

func (b *Builder) GqlServer(ctx context.Context) (*http.Server, error) {
	const timeout = time.Millisecond * 25

	router := b.httpRouter()

	router.HandleFunc(readinessEndpoint, b.healthcheck.handler)

	//nolint:exhaustruct
	server := http.Server{
		Addr:              b.config.HTTPAddr(),
		ReadHeaderTimeout: timeout,
		Handler:           router,
		ErrorLog:          log.New(zerolog.Nop(), "", 0),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	b.shutdown.add(func(ctx context.Context) error {
		return errors.Wrap(server.Shutdown(ctx), "shutdown graphql server")
	})

	return &server, nil
}

func (b *Builder) gqlHTTPRouter() *mux.Router {
	if b.http.router != nil {
		return b.http.router
	}

	b.http.router = mux.NewRouter()

	return b.http.router
}

func (b *Builder) SetGqlHandlers() error {
	router := b.httpRouter()
	router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{
			"content-type",
			"Access-Control-Request-Headers",
			"X-Requested-With",
		}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET"}),
	))

	handler := graphqlv1.NewGqlHandler()
	router.Handle("/query", handler)
	router.Handle("/playground", playground.Handler("Playground", "/query"))

	return nil
}
