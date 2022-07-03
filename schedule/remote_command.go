package schedule

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var RemoteCommand = &cli.Command{
	Name:        "remote",
	Aliases:     []string{"r"},
	Usage:       "send data to remote server as client",
	Description: "Currently, only support http-post",
	BashComplete: func(c *cli.Context) {
		for _, subCommand := range c.Command.Subcommands {
			fmt.Fprintf(c.App.Writer, "%s\n", subCommand.Name)
		}
	},
	Before: func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "doo begin ----- brace for impact\n")
		return nil
	},
	After: func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "doo end ----- did we lose anyone?\n")
		return nil
	},
	Action: func(c *cli.Context) error {
		fmt.Println(c.Command.FullName())
		fmt.Println(c.Command.HasName("wop"))
		fmt.Println(c.Command.Names())
		fmt.Println(c.Command.VisibleFlags())
		fmt.Printf("args : %t \n", c.Args())
		fmt.Fprintf(c.App.Writer, "dodododododoodododddooooododododooo\n")
		if c.Bool("forever") {
			c.Command.Run(c)
		}
		return nil
	},
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Command wrong, please use 'remote help' to check ... \n")
		return err
	},
}
