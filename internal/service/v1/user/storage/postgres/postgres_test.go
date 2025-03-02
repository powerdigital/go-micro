package postgres

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/powerdigital/go-micro/internal/config"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	conf, _ := config.Load()

	dsn := conf.Postgres.DSN()
	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	db.MustExec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT NOT NULL,
		age INT NOT NULL
	)`) // Ensure table exists for tests

	t.Cleanup(func() {
		db.MustExec(`DROP TABLE IF EXISTS users`)
		db.Close()
	})

	return db
}

func TestUserRepo(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	// Test CreateUser
	user := storage.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "123-456-7890",
		Age:   30,
	}

	userID, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	// Test GetUser
	fetchedUser, err := repo.GetUser(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, user.Name, fetchedUser.Name)
	require.Equal(t, user.Email, fetchedUser.Email)
	require.Equal(t, user.Phone, fetchedUser.Phone)
	require.Equal(t, user.Age, fetchedUser.Age)

	// Test GetUsers
	limit := 10
	users, err := repo.GetUsers(ctx, rune(limit))
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, user.Name, users[0].Name)

	// Test UpdateUser
	updatedUser := *fetchedUser
	updatedUser.Name = "Jane Doe"
	err = repo.UpdateUser(ctx, updatedUser)
	require.NoError(t, err)

	fetchedUpdatedUser, err := repo.GetUser(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "Jane Doe", fetchedUpdatedUser.Name)

	// Test DeleteUser
	err = repo.DeleteUser(ctx, userID)
	require.NoError(t, err)

	_, err = repo.GetUser(ctx, userID)
	require.Error(t, err)
}
