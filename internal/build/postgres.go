//nolint:revive
package build

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const maxPostgresTimeout = 10 * time.Second

func NewPostgresConnection(ctx context.Context, dsn string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, maxPostgresTimeout)
	defer cancel()

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
