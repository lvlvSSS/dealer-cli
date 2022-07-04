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
	},
)

var repeatFlag = &cli.IntFlag{
	Name:  "repeat",
	Usage: "specify the repeat times that the cron job do",
	Value: 1,
}

var repeatYamlFlag = altsrc.NewIntFlag(
	&cli.IntFlag{
		Name:  "dealer.schedule.repeat",
		Usage: "specify the repeat times that the cron job do, same as the flag 'repeat'",
	},
)

var timesFlag = &cli.IntFlag{
	Name:  "times",
	Usage: "specify the times while the cron job could do",
	Value: -1,
}
var timesYamlFlag = altsrc.NewIntFlag(
	&cli.IntFlag{
		Name:  "dealer.schedule.times",
		Usage: "specify the times while the cron job could do, same as the flag 'times' ",
	},
)

var durationFlag = &cli.Int64Flag{
	Name:  "duration",
	Usage: "specify the duration while the cron job could do, unit is second.",
	Value: -1,
}
var durationYamlFlag = altsrc.NewInt64Flag(
	&cli.Int64Flag{
		Name:  "dealer.schedule.duration",
		Usage: "specify the duration while the cron job could do, unit is second, same as the flag 'duration'",
	},
)
