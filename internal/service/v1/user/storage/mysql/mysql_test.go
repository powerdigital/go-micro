package mysql

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/powerdigital/go-micro/internal/config"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

var db *sql.DB
var repo storage.UserRepo

func TestMain(m *testing.M) {
	conf, _ := config.Load()

	var err error
	db, err = sql.Open("mysql", conf.MySQL.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		age INT NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo = NewUserRepo(db)

	m.Run()

	_, _ = db.Exec(`DROP TABLE IF EXISTS users;`)
}

func setupTestData() {
	_, _ = db.Exec(`TRUNCATE TABLE users;`)
}

func TestCreateUser(t *testing.T) {
	setupTestData()

	user := storage.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234567890",
		Age:   30,
	}

	id, err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	assert.NotZero(t, id)

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ?`, id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetUser(t *testing.T) {
	setupTestData()

	res, err := db.Exec(`INSERT INTO users (name, email, phone, age) VALUES (?, ?, ?, ?)`,
		"Jane Doe", "jane.doe@example.com", "0987654321", 25)
	require.NoError(t, err)

	id, err := res.LastInsertId()
	require.NoError(t, err)

	user, err := repo.GetUser(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Jane Doe", user.Name)
	assert.Equal(t, "jane.doe@example.com", user.Email)
	assert.Equal(t, "0987654321", user.Phone)
	assert.Equal(t, 25, user.Age)
}

func TestGetUsers(t *testing.T) {
	setupTestData()

	_, _ = db.Exec(`INSERT INTO users (name, email, phone, age) VALUES 
		('Alice', 'alice@example.com', '1111111111', 28),
		('Bob', 'bob@example.com', '2222222222', 35)`)

	users, err := repo.GetUsers(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUpdateUser(t *testing.T) {
	setupTestData()

	res, err := db.Exec(`INSERT INTO users (name, email, phone, age) VALUES (?, ?, ?, ?)`,
		"Charlie", "charlie@example.com", "3333333333", 40)
	require.NoError(t, err)

	id, err := res.LastInsertId()
	require.NoError(t, err)

	err = repo.UpdateUser(context.Background(), storage.User{
		ID:    id,
		Name:  "Charlie Updated",
		Email: "updated@example.com",
		Phone: "4444444444",
		Age:   45,
	})
	require.NoError(t, err)

	var name, email, phone string
	var age int
	err = db.QueryRow(`SELECT name, email, phone, age FROM users WHERE id = ?`, id).Scan(&name, &email, &phone, &age)
	require.NoError(t, err)
	assert.Equal(t, "Charlie Updated", name)
	assert.Equal(t, "updated@example.com", email)
	assert.Equal(t, "4444444444", phone)
	assert.Equal(t, 45, age)
}

func TestDeleteUser(t *testing.T) {
	setupTestData()

	res, err := db.Exec(`INSERT INTO users (name, email, phone, age) VALUES (?, ?, ?, ?)`,
		"Dave", "dave@example.com", "5555555555", 50)
	require.NoError(t, err)

	id, err := res.LastInsertId()
	require.NoError(t, err)

	err = repo.DeleteUser(context.Background(), id)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ?`, id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
