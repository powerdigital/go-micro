//nolint:revive
package cmd

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

func postgresCmd(ctx context.Context, conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "postgres",
		Short: "run db migrations for postgres",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(up(ctx, conf, postgres))
	command.AddCommand(down(ctx, conf, postgres))

	return command
}

func postgres(_ context.Context, conf config.Config) (*migrate.Migrate, error) {
	b := build.New(conf)

	//nolint:wrapcheck
	return b.PostgresMigration()
}
