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

func getMigrator() *migrate.Migrator {
	db := models.GetDB()
	return migrate.New(db, path)
}

// RunCreate create migration stub
func RunCreate(c *cli.Context) error {
	name := c.Args().First()
	migrator := getMigrator()
	migrator.SetName(name)
	return migrator.Create()
}
