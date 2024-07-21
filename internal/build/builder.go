package build

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/powerdigital/go-micro/internal/config"
	"github.com/powerdigital/go-micro/internal/service"
)

type Builder struct {
	config config.Config

	service service.Service

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
