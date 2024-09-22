package build

import (
	"io"

	sentry "github.com/archdx/zerolog-sentry"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Plugin func() (zerolog.LevelWriter, error)

func NewLogger(w io.Writer, plugins ...Plugin) zerolog.Logger {
	writers := []io.Writer{w}
	errors := make([]error, 0, len(plugins))

	for _, p := range plugins {
		if writer, err := p(); err != nil {
			errors = append(errors, err)
		} else {
			writers = append(writers, writer)
		}
	}

	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi)

	for _, err := range errors {
		logger.Err(err).Send()
	}

	return logger
}

func SentryWriter(dsn string) func() (zerolog.LevelWriter, error) {
	return func() (zerolog.LevelWriter, error) {
		if dsn == "" {
			return nil, errors.New("sentry dsn is not configured") //nolint:goerr113
		}

		w, err := sentry.New(dsn)
		if err != nil {
			return nil, errors.Wrapf(err, "init sentry connection")
		}

		return w, nil
	}
}
