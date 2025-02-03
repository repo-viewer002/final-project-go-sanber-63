package migrations

import (
	"embed"
)

//go:embed *.sql
var MigrationsDirectory embed.FS