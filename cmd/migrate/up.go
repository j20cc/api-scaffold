package migrate

import (
	"github.com/urfave/cli/v2"
)

// RunUp up migrations
func RunUp(c *cli.Context) error {
	name := c.Args().First()
	migrator.SetName(name)
	return migrator.Up()
}
