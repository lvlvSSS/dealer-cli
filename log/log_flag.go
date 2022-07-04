// log_flag doesn't mean that dealer-cli has log flag.
// log_flag is to analysis the mode and output flag to initialize the logger.
package log

import (
	"dealer-cli/utils/log"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

var ModeFlag = &cli.StringFlag{
	Name:        "mode",
	EnvVars:     []string{"MODE"},
	Usage:       "define the mode of the tool should work on .",
	DefaultText: "the value could be 'dev', 'debug', 'production' ",
}
var ModeYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "dealer.mode",
	Usage: "define the mode of the tool should work on . same as the flag 'mode' ",
})

var location, _ = os.Getwd()
var LogPathFlag = &cli.StringFlag{
	Name:  "log",
	Usage: "specify the log file path",
	Value: filepath.Join(location, "Logs/dealer-cli.log"),
}
var LogPathYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "dealer.log",
	Usage: "specify the log file path. same as the flag 'log' ",
})

func InitLog(c *cli.Context) error {
	var mode = ""
	if mode = c.String("mode"); len(strings.TrimSpace(mode)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - mode[%s]", mode))
	} else if mode = c.String("dealer.mode"); len(strings.TrimSpace(mode)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - dealer.mode[%s]", mode))
	} else {
		mode = "info"
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - mode forced to [%s]", mode))
	}
	var logPath = ""
	if LogPathFlag.IsSet() {
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - log[%s]", logPath))
		logPath = c.String("log")
	} else if logPath = c.String("dealer.log"); len(strings.TrimSpace(logPath)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - dealer.log[%s]", logPath))
	} else {
		logPath = c.String("log")
		log.Debug(fmt.Sprintf("[dealer_cli.log.InitLog] - log forced to [%s]", logPath))
	}

	level := zap.InfoLevel
	if strings.Compare(strings.ToLower(mode), "debug") == 0 ||
		strings.Compare(strings.ToLower(mode), "dev") == 0 {
		level = zap.DebugLevel
	}
	logPath, err := filepath.Abs(logPath)
	if err != nil {
		return err
	}
	log.Init(logPath, level)
	return nil
}
