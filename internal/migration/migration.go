package migration

import "embed"

//go:embed *
var FS embed.FS

const EmbedTypeIofs = "iofs"

const (
	MysqlSource    = "mysql"
	PostgresSource = "postgres"
)
