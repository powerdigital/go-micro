package entity

import "github.com/powerdigital/go-micro/internal/service/v1/user/storage"

type User struct {
	ID    int64  `json:"id,omitempty"`
	Name  string `json:"name" validate:"required,real_name"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,e164"`
	Age   int    `json:"age" validate:"required,min=18,max=65"`
}

type ErrResponse struct {
	Error string `json:"error"`
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
