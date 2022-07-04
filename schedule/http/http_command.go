package http

import (
	"dealer-cli/docs"
	LOG "dealer-cli/log"
	schedule_internal "dealer-cli/schedule/internal"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"regexp"
	"strings"
)

var httpFlags = []cli.Flag{
	headerFlag,
	headerYamlFlag,
}
var HttpCommand = &cli.Command{
	Name:  "http",
	Usage: "send http messages to server.",
	Flags: httpFlags,
	Action: func(c *cli.Context) error {
		if err := LOG.InitLog(c); err != nil {
			log.Warn(fmt.Sprintf("dealer_cli schedule http - log init errors[%s]", err))
			return err
		}
		if c.Bool("check") || c.Bool("dealer.check") {
			checkFlag(c)
			return nil
		}
		cronSchedule, err := checkCron(c)
		if err != nil {
			return err
		}

		// use FileRequestProducer
		producer := &FileRequestProducer{}
		httpClient := New()
		httpClient.After(producer.After)
		runner, err := checkRunner(c, producer, httpClient)
		if err != nil {
			return err
		}
		schedule_internal.DefaultCron.Schedule(cronSchedule, runner)
		schedule_internal.DefaultCron.Start()
		runner.Done()
		return nil
	},
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Command wrong, please use 'schedule http help' to check ... \n")
		return err
	},
	Before: altsrc.InitInputSourceWithContext(httpFlags, altsrc.NewYamlSourceFromFlagFunc(docs.APP_LOAD_YAML)),
}

// check the cron is valid or not
func checkCron(c *cli.Context) (cron.Schedule, error) {
	cronExpr := ""
	if cronExpr = c.String("cron"); len(strings.TrimSpace(cronExpr)) != 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - cron[%s]", cronExpr))
	} else if cronExpr = c.String("dealer.schedule.cron"); len(strings.TrimSpace(cronExpr)) != 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - dealer.schedule.cron[%s]", cronExpr))
	} else {
		return nil, errors.New("dealer_cli schedule http - cron is empty")
	}
	cron, err := cron.ParseStandard(cronExpr)
	if err != nil {
		log.Error(fmt.Sprintf("dealer_cli schedule http - cron[%s] is invalid", cronExpr))
		return nil, err
	}
	return cron, nil
}

func checkFlag(c *cli.Context) {
	log.Warn("dealer-cli schedule http - check flags begins ...")

	headers := c.StringSlice("header")
	log.Info(fmt.Sprintf("header : %v, length[%d]", headers, len(headers)))

	headers = c.StringSlice("dealer.schedule.http.header")
	log.Info(fmt.Sprintf("header yaml : %v, length[%d]", headers, len(headers)))
}

var headerFilePathRegexExpr = `"([\s\S]+?)={1}?([\s\S]+?)"`

func analysisHeaderFilePath(value string) (map[string]string, error) {
	headerFilePathRegex, err := regexp.Compile(headerFilePathRegexExpr)
	if err != nil {
		return nil, err
	}

	submatches := headerFilePathRegex.FindAllStringSubmatch(value, -1)
	re := make(map[string]string, len(submatches))
	for _, submatch := range submatches {
		if len(submatch) != 3 {
			return nil, errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", value, headerFilePathRegexExpr))
		}
		re[submatch[1]] = submatch[2]
	}
	return re, nil
}

func analysisHeader(header string) (string, string, error) {
	headerFilePathRegex, err := regexp.Compile(headerFilePathRegexExpr)
	if err != nil {
		return "", "", err
	}
	submatches := headerFilePathRegex.FindAllStringSubmatch(header, -1)
	if len(submatches) != 1 {
		return "", "", errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", header, headerFilePathRegexExpr))
	}
	if len(submatches[0]) != 3 {
		return "", "", errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", header, headerFilePathRegexExpr))
	}
	return submatches[0][1], submatches[0][2], nil
}
