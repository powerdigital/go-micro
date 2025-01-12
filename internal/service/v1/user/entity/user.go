package entity

import "github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql/model"

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
}

func (u User) EntityToModel() model.User {
	return model.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   u.Age,
	}
}

func ModelToEntity(u *model.User) User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   u.Age,
	}
}
