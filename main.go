package main

import (
	"context"
	"os"

	"github.com/powerdigital/go-micro/cmd"
	"github.com/powerdigital/go-micro/internal/build"
	"github.com/powerdigital/go-micro/internal/config"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logLevel, err := conf.LogLevel()
	if err != nil {
		panic(err)
	}

	sentryWriter := build.SentryWriter(conf.Monitoring.SentryDSN)
	log := build.NewLogger(os.Stdout, sentryWriter).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx := log.WithContext(context.Background())

	log.Info().Msg("the application is launching")

	exitCode := 0

	err = cmd.Run(ctx, conf)
	if err != nil {
		log.Err(err).Send()

		exitCode = 1
	}

	os.Exit(exitCode)
}
