package build

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/powerdigital/go-micro/internal/rest"
)

func (b *Builder) HTTPServer(ctx context.Context) (*http.Server, error) {
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
		return errors.Wrap(server.Shutdown(ctx), "shutdown http server")
	})

	return &server, nil
}

func (b *Builder) httpRouter() *mux.Router {
	if b.http.router != nil {
		return b.http.router
	}

	b.http.router = mux.NewRouter()

	return b.http.router
}

func (b *Builder) SetHTTPHandlers() error {
	router := b.httpRouter()

	service, err := b.GreetingService()
	if err != nil {
		return errors.Wrap(err, "get service")
	}

	handler := rest.NewHandler(service)

	router.HandleFunc("/", handler.GetHello)

	return nil
}
