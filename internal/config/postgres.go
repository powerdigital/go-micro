package config

import (
	"fmt"
	"net"
	"strconv"
)

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true" default:"localhost"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true" default:"5432"`
	Username string `envconfig:"POSTGRES_USER" required:"true" default:"micro"`
	Password string `envconfig:"POSTGRES_PASS" required:"true" default:"secret"`
	Database string `envconfig:"POSTGRES_BASE" required:"true" default:"micro"`
}

func (pg *Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pg.Host,
		pg.Port,
		pg.Username,
		pg.Password,
		pg.Database,
	)
}

func (pg *Postgres) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		pg.Username,
		pg.Password,
		net.JoinHostPort(pg.Host, strconv.Itoa(pg.Port)),
		pg.Database,
	)
}
