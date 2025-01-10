package config

type MySQL struct {
	DSN string `envconfig:"MYSQL_DSN" default:"micro:secret@(localhost:3306)/micro"`
}
