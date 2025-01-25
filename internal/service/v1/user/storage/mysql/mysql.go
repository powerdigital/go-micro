package mysql

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"

	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) storage.UserRepo {
	return &userRepo{db: db}
}

func (repo *userRepo) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	query := `INSERT INTO users (name, email, phone, age) VALUES (?, ?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query, user.Name, user.Email, user.Phone, user.Age)
	if err != nil {
		return 0, errors.Wrap(err, "error inserting user")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "error getting insert id")
	}

	return id, nil
}

func (repo *userRepo) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	query := `SELECT id, name, email, phone, age FROM users WHERE id = ?`

	var user storage.User

	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, errors.Wrap(err, "error getting user")
	}

	return &user, nil
}

func (repo *userRepo) GetUsers(ctx context.Context) ([]storage.User, error) {
	query := `SELECT id, name, email, phone, age FROM users`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "error fetching users")
	}
	defer rows.Close()

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
	query := `UPDATE users SET name = ?, email = ?, phone = ?, age = ? WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, query, user.Name, user.Email, user.Phone, user.Age, user.ID)

	return errors.Wrap(err, "error updating user")
}

func (repo *userRepo) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, query, userID)

	return errors.Wrap(err, "error deleting user")
}
