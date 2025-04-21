//nolint:revive
package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/multierr"

	"github.com/powerdigital/go-micro/internal/migration"
)

func (b *Builder) PostgresMigration() (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(migration.FS, migration.PostgresSource)
	if err != nil {
		return nil, errors.Wrap(err, "embed postgres migrations")
	}

	migrator, err := migrate.NewWithSourceInstance(migration.EmbedTypeIofs, sourceDriver, b.config.Postgres.URL())
	if err != nil {
		return nil, errors.Wrap(err, "apply postgres migrations")
	}

	b.shutdown.add(func(ctx context.Context) error {
		return multierr.Append(migrator.Close())
	})

	return migrator, nil
}
