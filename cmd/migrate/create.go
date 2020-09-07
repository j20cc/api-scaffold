package migrate

import (
	"gvue-scaffold/app/models"
	"gvue-scaffold/internal/migrate"

	"github.com/urfave/cli/v2"
)

var (
	path     = "resources/migrations"
	tbname   = "migrations"
	migrator *migrate.Migrator
)

func init() {
	db := models.GetDB()
	migrator = migrate.New(db, path)
}

// RunCreate create migration stub
func RunCreate(c *cli.Context) error {
	name := c.Args().First()
	migrator.SetName(name)
	return migrator.Create()
}
