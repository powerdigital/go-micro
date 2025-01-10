package converter

import (
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql/model"
)

func ModelToEntity(u *model.User) entity.User {
	return entity.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   u.Age,
	}
}
