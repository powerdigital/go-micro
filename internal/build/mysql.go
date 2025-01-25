//nolint:revive
package build

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const maxMySQLTimeout = 10 * time.Second

func NewMySQLConnection(ctx context.Context, dsn string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, maxMySQLTimeout)
	defer cancel()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
