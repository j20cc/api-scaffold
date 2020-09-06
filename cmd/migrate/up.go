package migrate

import (
	"gvue-scaffold/app/models"
	"gvue-scaffold/internal/migrate"

	"github.com/urfave/cli/v2"
)

// RunUp up migrations
// migration id migration bath
func RunUp(c *cli.Context) error {
	db := models.GetDB()
	name := c.Args().First()
	migrator := migrate.New(db, path, name)
	return migrator.Up()
}
