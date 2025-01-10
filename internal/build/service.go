package build

import (
	helloservice "github.com/powerdigital/go-micro/internal/service/v1/greeting"
	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql"
)

func (b *Builder) GreetingService() (helloservice.HelloSrv, error) {
	if b.greetingService == nil {
		b.greetingService = helloservice.NewHelloService()
	}

	return b.greetingService, nil
}

func (b *Builder) UserService() (userservice.UserSrv, error) {
	if b.userService == nil {
		db, err := NewMySQLConnection(b.config)
		if err != nil {
			return nil, err
		}

		repo := mysql.NewUserRepo(db)

		b.userService = userservice.NewUserService(repo)
	}

	return b.userService, nil
}
