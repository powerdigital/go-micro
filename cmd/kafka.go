package cmd

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

type consumerStorage func(ctx context.Context) (*build.Consumer, error)

func kafkaServer(ctx context.Context, conf config.Config) *cobra.Command {
	builder := build.New(conf)

	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "start kafka consumer",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			zerolog.Ctx(ctx).Info().Msg("start " + cmd.Short)
		},
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			zerolog.Ctx(ctx).Info().Msg("stop " + cmd.Short)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "all",
			Short: "common consumer for reading all topics",
			RunE: runConsumers(ctx, []consumerStorage{
				builder.UserCreatingConsumer,
			}, builder),
		},
	)

	return cmd
}

//nolint:cyclop
func runConsumers(
	ctx context.Context,
	consumers []consumerStorage,
	builder *build.Builder,
) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		stop := builder.ChannelShutdown(ctx)

		srv, err := builder.HTTPServer(ctx)
		if err != nil {
			return errors.Wrap(err, "build http server")
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				zerolog.Ctx(ctx).Err(err).Msg("run http server")
			}
		}()

		wg := sync.WaitGroup{}

		localCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		for _, consumerBlank := range consumers {
			consumer, err := consumerBlank(localCtx)
			if err != nil {
				cancel()

				return err
			}

			wg.Add(1)

			go func(consumer *build.Consumer) {
				defer wg.Done()

				logger := zerolog.Ctx(localCtx)
				ctx := logger.With().Str("topic", consumer.Topic).Logger().WithContext(localCtx)

				for {
					select {
					case <-ctx.Done():
						return
					case <-stop:
						return
					default:
						if err := consumer.Consume(ctx); err != nil && !errors.Is(err, sarama.ErrClosedConsumerGroup) {
							zerolog.Ctx(ctx).Err(err).Msg("run consumer")
						}
					}

					time.Sleep(time.Second)
				}
			}(consumer)
		}

		wg.Wait()

		return nil
	}
}
