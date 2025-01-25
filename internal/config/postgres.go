package config

import "fmt"

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true" default:"localhost"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true" default:"5432"`
	Username string `envconfig:"POSTGRES_USER" required:"true" default:"micro"`
	Password string `envconfig:"POSTGRES_PASS" required:"true" default:"secret"`
	Database string `envconfig:"POSTGRES_NAME" required:"true" default:"micro"`
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
