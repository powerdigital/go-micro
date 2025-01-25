package cmd

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

func gqlServer(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "graphql",
		Short: "start graphql server",
		RunE: func(cmd *cobra.Command, args []string) error { //nolint:revive
			builder := build.New(conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			srv, err := builder.GqlServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build graphql server")
			}

			err = builder.SetGqlHandlers()
			if err != nil {
				return errors.Wrap(err, "set graphql handlers")
			}

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				zerolog.Ctx(ctx).Err(errors.WithStack(err)).Msg("run graphql server")
			}

			<-ctx.Done()

			return nil
		},
	}
}
