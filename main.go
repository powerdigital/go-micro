package main

import (
	"context"
	"os"

	"github.com/powerdigital/go-micro/cmd"
	"github.com/powerdigital/go-micro/internal/config"
	"github.com/rs/zerolog"
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

	log := zerolog.New(os.Stdout).Level(logLevel)

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
