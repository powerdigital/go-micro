package build

import (
	"context"

	helloservice "github.com/powerdigital/go-micro/internal/service/v1/greeting"
	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/postgres"
)

func (b *Builder) GreetingService() (helloservice.HelloSrv, error) {
	if b.greetingService == nil {
		b.greetingService = helloservice.NewHelloService()
	}

	return b.greetingService, nil
}

func (b *Builder) UserService(ctx context.Context) (userservice.UserSrv, error) {
	if b.userService != nil {
		return b.userService, nil
	}

	var repo storage.UserRepo

	if b.config.App.Storage == "postgres" {
		db, err := NewPostgresConnection(ctx, b.config.Postgres.DSN())
		if err != nil {
			return nil, err
		}

		repo = postgres.NewUserRepo(db)
	} else {
		db, err := NewMySQLConnection(ctx, b.config.MySQL.DSN())
		if err != nil {
			return nil, err
		}

		repo = mysql.NewUserRepo(db)
	}

	b.userService = userservice.NewUserService(repo)

	return b.userService, nil
}
