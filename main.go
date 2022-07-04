package main

import (
	"dealer-cli/docs"
	dealer_file "dealer-cli/file"
	"dealer-cli/log"
	"dealer-cli/schedule"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"time"
)

var checkFlag = &cli.BoolFlag{
	Name:  "check",
	Usage: "used to check flags' value",
	Value: false,
}

var checkYamlFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:  "dealer.check",
	Usage: "used to check flags' value. same as the flag 'check' ",
})

var appFlags = []cli.Flag{
	log.ModeFlag,
	log.ModeYamlFlag,
	log.LogPathFlag,
	log.LogPathYamlFlag,
	checkFlag,
	checkYamlFlag,
}

func main() {
	app := &cli.App{
		Name:     "dealer-cli",
		Version:  "v0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Nelson Lv",
			},
		},
		Usage: "a simple cli automatic tool.",
		// global error handler
		ExitErrHandler: func(c *cli.Context, err error) {
			fmt.Printf("ERROR: Command[%s] - error[%s] \n", c.Command.Name, err.Error())
		},
		Flags: append(appFlags, docs.LoadFlag),
		Before: func(c *cli.Context) error {
			// if the 'load-yaml' is not specified, then use the default value.
			// Must Set the 'load-yaml', otherwise the 'load-yaml' can't affect.
			if !c.IsSet(docs.APP_LOAD_YAML) {
				c.Set(docs.APP_LOAD_YAML, docs.LoadFlag.Value)
			}
			var before = altsrc.InitInputSourceWithContext(appFlags, altsrc.NewYamlSourceFromFlagFunc(docs.APP_LOAD_YAML))
			return before(c)
		},
		Commands: []*cli.Command{
			schedule.ScheduleCommand,
			dealer_file.FileCommand,
		},
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,

		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "[%q] not supported now !\n", command)
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}
			fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("App[%s] occurs error[%s] ... ", app.Name, err.Error())
	}
}
