package v1

import (
	servicev1 "github.com/powerdigital/go-micro/internal/service/v1/greeting"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service servicev1.HelloSrv
}

func NewResolver(service servicev1.HelloSrv) *Resolver {
	return &Resolver{
		service: service,
	}
}
