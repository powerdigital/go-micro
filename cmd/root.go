package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/config"
)

func Run(ctx context.Context, conf config.Config) error {
	root := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error { //nolint:revive
			return cmd.Usage()
		},
	}

	root.AddCommand(
		restServer(ctx, conf),
		grpcServer(ctx, conf),
		gqlServer(ctx, conf),
		kafkaServer(ctx, conf),
		migrateCmd(ctx, conf),
	)

	return errors.Wrap(root.ExecuteContext(ctx), "run application")
}
