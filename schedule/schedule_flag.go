package schedule

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"path/filepath"
)

var location, _ = os.Getwd()
var cronFlag = &cli.StringSliceFlag{
	Name:     "cron",
	Usage:    "schedule the job by cron expr",
	FilePath: filepath.Join(location, "./schedule.cron.dealer"),
}

var cronYamlFlag = altsrc.NewStringSliceFlag(
	&cli.StringSliceFlag{
		Name:  "dealer.schedule.cron",
		Usage: "schedule the job by cron expr, same as the flag 'cron'",
	})

var repeatFlag = &cli.IntFlag{
	Name:  "repeat",
	Usage: "specify the repeat times that the cron job do",
	Value: 1,
}

var repeatYamlFlag = altsrc.NewIntFlag(
	&cli.IntFlag{
		Name:  "dealer.schedule.repeat",
		Usage: "specify the repeat times that the cron job do, same as the flag 'repeat'",
		Value: 1,
	})
