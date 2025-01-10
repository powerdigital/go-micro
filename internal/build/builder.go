package build

import (
	"github.com/gorilla/mux"

	"github.com/powerdigital/go-micro/internal/config"
	greetingservice "github.com/powerdigital/go-micro/internal/service/v1/greeting"
	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
)

type Builder struct {
	config config.Config

	greetingService greetingservice.HelloSrv
	userService     userservice.UserSrv

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
