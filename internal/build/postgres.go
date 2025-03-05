//nolint:revive,dupl
package build

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const maxPostgresTimeout = 10 * time.Second

func (b *Builder) newPostgresConnection(ctx context.Context, dsn string) (*sqlx.DB, error) {
	_, cancel := context.WithTimeout(ctx, maxPostgresTimeout)
	defer cancel()

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection: %w", err)
	}

	b.shutdown.add(func(_ context.Context) error {
		if err = db.Close(); err != nil {
			return errors.Wrap(err, "close postgres db connection")
		}

		return nil
	})

	b.healthcheck.add(func(_ context.Context) error {
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, "ping postgres db")
		}

		return nil
	})

	return db, nil
}
