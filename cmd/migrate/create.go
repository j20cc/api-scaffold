package migrate

import (
	"gvue-scaffold/app/models"
	"gvue-scaffold/internal/migrate"

	"github.com/urfave/cli/v2"
)

var (
	path   = "resources/migrations"
	tbname = "migrations"
)

// RunCreate create migration stub
func RunCreate(c *cli.Context) error {
	db := models.GetDB()
	name := c.Args().First()
	migrator := migrate.New(db, path, name)
	return migrator.Create()
}
