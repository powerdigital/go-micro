package postgres

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"

	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) storage.UserRepo {
	return &userRepo{db: db}
}

func (repo *userRepo) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	query := `INSERT INTO users (name, email, phone, age) VALUES (:name, :email, :phone, :age) RETURNING id`

	result, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"age":   user.Age,
	})
	if err != nil {
		return 0, errors.Wrap(err, "error inserting user")
	}

	defer func() {
		_ = result.Close()
		_ = result.Err()
	}()

	result.Next()

	if err := result.Scan(&user.ID); err != nil {
		return 0, fmt.Errorf("event retrieve last insert id: %w", err)
	}

	return user.ID, nil
}

func (repo *userRepo) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	query := `SELECT id, name, email, phone, age FROM users WHERE id = :id`

	var user storage.User

	result, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting user")
	}

	defer func() {
		_ = result.Close()
		_ = result.Err()
	}()

	if !result.Next() {
		return nil, storage.ErrNotFound
	}

	if err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Age); err != nil {
		return nil, errors.Wrap(err, "error getting user")
	}

	return &user, nil
}

func (repo *userRepo) GetUsers(ctx context.Context) ([]storage.User, error) {
	query := `SELECT id, name, email, phone, age FROM users`

	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{})
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "error fetching users")
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var users []storage.User

	for rows.Next() {
		var user storage.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Age); err != nil {
			return nil, errors.Wrap(err, "error scanning user")
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepo) UpdateUser(ctx context.Context, user storage.User) error {
	query := `UPDATE users SET name = :name, email = :email, phone = :phone, age = :age WHERE id = :id`

	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"age":   user.Age,
		"id":    user.ID,
	})

	return errors.Wrap(err, "error updating user")
}

func (repo *userRepo) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": userID,
	})

	return errors.Wrap(err, "error deleting user")
}
