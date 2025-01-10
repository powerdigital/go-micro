//nolint:revive
package build

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/powerdigital/go-micro/internal/config"
)

func NewMySQLConnection(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.MySQL.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
