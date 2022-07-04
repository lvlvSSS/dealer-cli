package schedule

import "github.com/urfave/cli/v2"

var cronFlag = &cli.StringSliceFlag{
	Name: "cron",
}
