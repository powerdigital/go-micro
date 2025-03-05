//nolint:revive,dupl
package build

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	_ "github.com/go-sql-driver/mysql"
)

const maxMySQLTimeout = 10 * time.Second

func (b *Builder) newMySQLConnection(ctx context.Context, dsn string) (*sql.DB, error) {
	_, cancel := context.WithTimeout(ctx, maxMySQLTimeout)
	defer cancel()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	b.shutdown.add(func(_ context.Context) error {
		if err = db.Close(); err != nil {
			return errors.Wrap(err, "close mysql db connection")
		}

		return nil
	})

	b.healthcheck.add(func(_ context.Context) error {
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, "ping mysql db")
		}

		return nil
	})

	return db, nil
}
