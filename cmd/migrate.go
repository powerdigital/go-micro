//nolint:revive
package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/config"
)

func migrateCmd(ctx context.Context, conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "migrate",
		Short:     "run db migrations",
		ValidArgs: []string{"postgres"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(
		mysqlCmd(ctx, conf),
		postgresCmd(ctx, conf),
	)

	return command
}

type migrationFn func(context.Context, config.Config) (*migrate.Migrate, error)

func up(ctx context.Context, conf config.Config, migrationFn migrationFn) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "execute all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrationFn(ctx, conf)
			if err != nil {
				return errors.Wrap(err, "create migration")
			}

			err = m.Up()
			if err != nil {
				if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
					return nil
				}

				return errors.Wrap(err, "execute migrations")
			}

			return nil
		},
	}
}

func down(ctx context.Context, conf config.Config, migrationFn migrationFn) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "rollback all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrationFn(ctx, conf)
			if err != nil {
				return errors.Wrap(err, "create migration")
			}

			err = m.Down()
			if err != nil {
				return errors.Wrap(err, "rollback migrations")
			}

			return nil
		},
	}
}
