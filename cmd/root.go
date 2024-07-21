package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/powerdigital/go-micro/internal/config"
	"github.com/spf13/cobra"
)

// nolint:revive
func Run(ctx context.Context, conf config.Config) error {
	root := &cobra.Command{ //nolint:exhaustruct
		RunE: func(cmd *cobra.Command, args []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	root.AddCommand(
		httpServer(ctx, conf),
	)

	return errors.Wrap(root.ExecuteContext(ctx), "run application")
}
