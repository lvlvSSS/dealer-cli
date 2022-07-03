// log_flag doesn't mean that dealer-cli has log flag.
// log_flag is to analysis the mode and output flag to initialize the logger.
package log

import (
	"dealer-cli/utils/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

var ModeFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:        "mode",
	EnvVars:     []string{"MODE"},
	Usage:       "define the mode of the tool should work on .",
	DefaultText: "the value could be 'dev', 'debug', 'production' ",
	Value:       "debug",
})

var location, _ = os.Getwd()
var LogPathFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "log",
	Usage: "specify the log file path",
	Value: filepath.Join(location, "dealer-cli.log"),
})

func InitLog(c *cli.Context) {
	mode := c.String("mode")
	output := c.String("log")
	level := zap.InfoLevel
	if strings.Compare(strings.ToLower(mode), "debug") == 0 ||
		strings.Compare(strings.ToLower(mode), "dev") == 0 {
		level = zap.DebugLevel
	}
	log.Init(output, level)
}
