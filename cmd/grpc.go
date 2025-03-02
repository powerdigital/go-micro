package cmd

import (
	"context"
	"net"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

func grpcServer(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "run grpc server",
		RunE: func(_ *cobra.Command, _ []string) error {
			builder := build.New(conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			grpcSrv, err := builder.GRPCServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build grpc server")
			}

			reflection.Register(grpcSrv)

			httpSrv, err := builder.HTTPServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build http server")
			}

			go func() {
				if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					zerolog.Ctx(ctx).Err(err).Msg("run http server")
				}
			}()

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			listener, err := net.Listen(conf.GRPCNetworkType(), conf.GRPCAddress())
			if err != nil {
				zerolog.Ctx(ctx).Err(errors.Wrap(err, "start network listener")).Send()
			}

			if err = grpcSrv.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
				zerolog.Ctx(ctx).Err(errors.Wrap(err, "run grpc server")).Send()
			}

			<-ctx.Done()

			return nil
		},
	}
}
