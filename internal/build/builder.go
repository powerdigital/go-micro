package build

import (
	"github.com/gorilla/mux"

	"github.com/powerdigital/go-micro/internal/config"
	servicev1 "github.com/powerdigital/go-micro/internal/service/v1/greeting"
)

type Builder struct {
	config config.Config

	service servicev1.HelloService

	shutdown    shutdown
	healthcheck healthcheck

	http struct {
		router *mux.Router
	}
}

func New(conf config.Config) *Builder {
	b := Builder{config: conf} //nolint:exhaustruct

	return &b
}
