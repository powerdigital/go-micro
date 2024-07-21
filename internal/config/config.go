package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	App  App
	HTTP HTTP
}

type App struct {
	ENV      string `default:"local" envconfig:"APP_ENV"`
	Name     string `default:"app"   envconfig:"APP_NAME"`
	LogLevel string `default:"debug" envconfig:"LOG_LEVEL"`
}

type HTTP struct {
	Port    int32    `default:"8080" envconfig:"HTTP_PORT"`
	Schemes []string `default:"http" envconfig:"HTTP_SCHEMES"`
}

func Load() (Config, error) {
	cnf := Config{} //nolint:exhaustruct

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cnf, errors.Wrap(err, "read .env file")
	}

	if err := envconfig.Process("", &cnf); err != nil {
		return cnf, errors.Wrap(err, "read environment")
	}

	return cnf, nil
}

func (c *Config) LogLevel() (zerolog.Level, error) {
	lvl, err := zerolog.ParseLevel(c.App.LogLevel)
	if err != nil {
		return 0, errors.Wrapf(err, "loading log level from config value %q", c.App.LogLevel)
	}

	return lvl, nil
}

func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTP.Port)
}