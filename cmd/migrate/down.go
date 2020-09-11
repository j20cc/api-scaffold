package migrate

import (
	"github.com/urfave/cli/v2"
)

// RunDown rollback migrations
func RunDown(c *cli.Context) error {
	name := c.Args().First()
	migrator := getMigrator()
	migrator.SetName(name)
	return migrator.Down()
}
