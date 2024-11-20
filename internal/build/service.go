package build

import (
	v1 "github.com/powerdigital/go-micro/internal/service/v1"
)

func (b *Builder) GreetingService() (v1.GreetingService, error) {
	if b.service == nil {
		b.service = v1.NewService()
	}

	return b.service, nil
}
