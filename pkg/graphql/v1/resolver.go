package v1

import v1 "github.com/powerdigital/go-micro/internal/service/v1"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service v1.GreetingService
}

func NewResolver(service v1.GreetingService) *Resolver {
	return &Resolver{
		service: service,
	}
}
