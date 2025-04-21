package config

import "fmt"

type MySQL struct {
	Host     string `envconfig:"MYSQL_HOST" default:"localhost"`
	Port     int    `envconfig:"MYSQL_PORT" default:"3306"`
	Username string `envconfig:"MYSQL_USER" default:"micro"`
	Password string `envconfig:"MYSQL_PASS" default:"secret"`
	Database string `envconfig:"MYSQL_BASE" default:"micro"`
}

func (my *MySQL) DSN() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s",
		my.Username,
		my.Password,
		my.Host,
		my.Port,
		my.Database,
	)
}

func (my *MySQL) URL() string {
	return "mysql://" + my.DSN()
}
