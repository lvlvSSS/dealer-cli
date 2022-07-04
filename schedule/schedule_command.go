package schedule

import (
	"dealer-cli/docs"
	"dealer-cli/schedule/http"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var scheduleFlags = []cli.Flag{
	cronFlag,
	cronYamlFlag,
	repeatFlag,
	repeatYamlFlag,
	timesFlag,
	timesYamlFlag,
	durationFlag,
	durationYamlFlag,
}

var ScheduleCommand = &cli.Command{
	Name:    "schedule",
	Aliases: []string{"s"},
	Usage:   "schedule job by using cron, schedule command always has the subcommand, to do some specified action",
	Flags:   scheduleFlags,
	BashComplete: func(c *cli.Context) {
		for _, subCommand := range c.Command.Subcommands {
			fmt.Fprintf(c.App.Writer, "%s\n", subCommand.Name)
		}
	},
	Before: altsrc.InitInputSourceWithContext(scheduleFlags, altsrc.NewYamlSourceFromFlagFunc(docs.APP_LOAD_YAML)),
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Command wrong, please use 'schedule help' to check ... \n")
		return err
	},
	Subcommands: []*cli.Command{
		http.HttpCommand,
	},
}
