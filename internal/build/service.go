package build

import "github.com/powerdigital/go-micro/internal/service"

func (b *Builder) Service() (service.Service, error) {
	if b.service == nil {
		b.service = service.NewService()
	}

	return b.service, nil
}
