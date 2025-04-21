//nolint:revive
package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/powerdigital/go-micro/internal/migration"
)

func (b *Builder) MysqlMigration() (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(migration.FS, migration.MysqlSource)
	if err != nil {
		return nil, errors.Wrap(err, "embed mysql migrations")
	}

	migrator, err := migrate.NewWithSourceInstance(migration.EmbedTypeIofs, sourceDriver, b.config.MySQL.URL())
	if err != nil {
		return nil, errors.Wrap(err, "apply mysql migrations")
	}

	b.shutdown.add(func(ctx context.Context) error {
		return errors.Join(migrator.Close())
	})

	return migrator, nil
}
