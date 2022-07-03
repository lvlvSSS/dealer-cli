package file

import (
	"dealer-cli/file/extract"
	"fmt"
	"github.com/urfave/cli/v2"
)

var FileCommand = &cli.Command{
	Name:        "file",
	Usage:       "file command always has the subcommand, to do some specified action",
	Description: "Some prep work for file operations",
	BashComplete: func(c *cli.Context) {
		for _, subCommand := range c.Command.Subcommands {
			fmt.Fprintf(c.App.Writer, "%s\n", subCommand.Name)
		}
	},
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Command wrong, please use 'file help' to check ... \n")
		return err
	},
	Subcommands: []*cli.Command{
		extract.ExtractCommand,
	},
}
