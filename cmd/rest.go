//nolint:dupl
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

func restServer(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "start rest server",
		RunE: func(cmd *cobra.Command, args []string) error { //nolint:revive
			builder := build.New(conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			srv, err := builder.HTTPServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build http server")
			}

			err = builder.SetHTTPHandlers(ctx)
			if err != nil {
				return errors.Wrap(err, "set rest handlers")
			}

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				zerolog.Ctx(ctx).Err(errors.WithStack(err)).Msg("run http server")
			}

			<-ctx.Done()

			return nil
		},
	}
}
