package build

import (
	servicev1 "github.com/powerdigital/go-micro/internal/service/v1/greeting"
)

func (b *Builder) GreetingService() (servicev1.HelloService, error) {
	if b.service == nil {
		b.service = servicev1.NewService()
	}

	return b.service, nil
}
