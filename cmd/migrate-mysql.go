//nolint:revive
package cmd

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

func mysqlCmd(ctx context.Context, conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "mysql",
		Short: "run db migrations for mysql",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(up(ctx, conf, mysql))
	command.AddCommand(down(ctx, conf, mysql))

	return command
}

func mysql(_ context.Context, conf config.Config) (*migrate.Migrate, error) {
	b := build.New(conf)

	//nolint:wrapcheck
	return b.MysqlMigration()
}
