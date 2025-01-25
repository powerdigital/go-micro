package entity

import "github.com/powerdigital/go-micro/internal/service/v1/user/storage"

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
}

func (u User) EntityToModel() storage.User {
	return storage.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   u.Age,
	}
}

func ModelToEntity(u *storage.User) User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   u.Age,
	}
}
