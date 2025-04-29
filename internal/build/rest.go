package build

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/powerdigital/go-micro/internal/middleware"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	restv1 "github.com/powerdigital/go-micro/internal/transport/rest/v1"
)

const (
	requestTimeout    = 5 * time.Second
	readHeaderTimeout = 25 * time.Millisecond
)

func (b *Builder) httpRouter() *mux.Router {
	if b.http.router != nil {
		return b.http.router
	}

	b.http.router = mux.NewRouter()

	return b.http.router
}

func (b *Builder) HTTPServer(ctx context.Context) (*http.Server, error) {
	router := b.httpRouter()

	router.HandleFunc(readinessEndpoint, b.healthcheck.handler)

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.SecurityHeadersMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.TimeoutMiddleware(requestTimeout))

	//nolint:exhaustruct
	server := http.Server{
		Addr:              b.config.HTTPAddr(),
		ReadHeaderTimeout: readHeaderTimeout,
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

func (b *Builder) SetHTTPHandlers(ctx context.Context) error {
	router := b.httpRouter()

	service, err := b.UserService(ctx)
	if err != nil {
		return errors.Wrap(err, "get UserService")
	}

	validate := validator.New()
	if err = entity.RegisterValidators(validate); err != nil {
		return errors.Wrap(err, "register entity validators")
	}

	handler := restv1.NewRESTHandler(service, validate)

	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	return nil
}
