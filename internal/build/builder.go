package build

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/powerdigital/go-micro/internal/config"
	v1 "github.com/powerdigital/go-micro/internal/service/v1"
)

type Builder struct {
	config config.Config

	service v1.GreetingService

	shutdown    shutdown
	healthcheck healthcheck

	http struct {
		router *mux.Router
		server *http.Server
	}
}

func New(conf config.Config) *Builder {
	b := Builder{config: conf} //nolint:exhaustruct

	return &b
}
