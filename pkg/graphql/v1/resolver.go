package v1

import userservice "github.com/powerdigital/go-micro/internal/service/v1/user"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service userservice.UserSrv
}

func NewResolver(service userservice.UserSrv) *Resolver {
	return &Resolver{
		service: service,
	}
}
