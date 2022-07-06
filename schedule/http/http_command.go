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
	"strings"
)

var httpFlags = []cli.Flag{
	headerFlag,
	headerYamlFlag,
	methodFlag,
	methodYamlFlag,
	bodyFlag,
	bodyYamlFlag,
	urlFlag,
	urlYamlFlag,
	doneLocationFlag,
	doneLocationYamlFlag,
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
		cronSchedules, err := checkCron(c)
		if err != nil {
			return err
		}

		// use FileRequestProducer
		producer := &FileRequestProducer{}
		producer.Init(c)

		httpClient := New()
		httpClient.After(producer.After)
		runner, err := checkRunner(c, producer, httpClient)
		if err != nil {
			return err
		}
		for _, cronSchedule := range cronSchedules {
			schedule_internal.DefaultCron.Schedule(cronSchedule, runner)
		}
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
func checkCron(c *cli.Context) ([]cron.Schedule, error) {
	cronExprs := make([]string, 0, 8)
	if cronExprs = c.StringSlice("cron"); len(cronExprs) != 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - cron[%v]", cronExprs))
	} else if cronExprs = c.StringSlice("dealer.schedule.cron"); len(cronExprs) != 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - dealer.schedule.cron[%v]", cronExprs))
	} else {
		return nil, errors.New("dealer_cli schedule http - cron is empty")
	}
	crons := make([]cron.Schedule, 0, len(cronExprs))
	for _, cronExpr := range cronExprs {
		cron, err := schedule_internal.DefaultParser.Parse(cronExpr)
		if err != nil {
			log.Error(fmt.Sprintf("dealer_cli schedule http - cron[%s] is invalid", cronExprs))
			return nil, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}

func checkFlag(c *cli.Context) {
	log.Warn("dealer-cli schedule http - check flags begins ...")

	cronExpr := c.StringSlice("cron")

	log.Info(fmt.Sprintf("cron : %v", cronExpr))

	cronExpr = c.StringSlice("dealer.schedule.cron")
	log.Info(fmt.Sprintf("cron yaml : %v", cronExpr))

	repeat := c.Int("repeat")
	log.Info(fmt.Sprintf("repeat : %v", repeat))

	repeat = c.Int("dealer.schedule.repeat")
	log.Info(fmt.Sprintf("repeat yaml : %v", repeat))

	times := c.Int("times")
	log.Info(fmt.Sprintf("times : %v", times))
	times = c.Int("dealer.schedule.times")
	log.Info(fmt.Sprintf("times yaml : %v", times))

	duration := c.Duration("duration")
	log.Info(fmt.Sprintf("duration : %v", duration))
	duration = c.Duration("dealer.schedule.duration")
	log.Info(fmt.Sprintf("duration yaml : %v", duration))

	headers := c.StringSlice("header")
	log.Info(fmt.Sprintf("header : %v, length[%d], set[%v]", headers, len(headers), c.IsSet("header")))

	headers = c.StringSlice("dealer.schedule.http.header")
	log.Info(fmt.Sprintf("header yaml : %v, length[%d]", headers, len(headers)))

	method := c.String("method")
	if len(strings.TrimSpace(method)) == 0 {
		method = "nil"
	}
	log.Info(fmt.Sprintf("method : %v", method))
	method = c.String("dealer.schedule.http.method")
	if len(strings.TrimSpace(method)) == 0 {
		method = "nil"
	}
	log.Info(fmt.Sprintf("method yaml : %v", method))

	body := c.String("body")
	if len(strings.TrimSpace(body)) == 0 {
		body = "nil"
	}
	log.Info(fmt.Sprintf("body : %v", body))
	body = c.String("dealer.schedule.http.body")
	if len(strings.TrimSpace(body)) == 0 {
		body = "nil"
	}
	log.Info(fmt.Sprintf("body yaml : %v", body))

	url := c.String("url")
	if len(strings.TrimSpace(url)) == 0 {
		url = "nil"
	}
	log.Info(fmt.Sprintf("url : %v", url))
	url = c.String("dealer.schedule.http.url")
	if len(strings.TrimSpace(url)) == 0 {
		url = "nil"
	}
	log.Info(fmt.Sprintf("url yaml : %v", url))

	doneLocation := c.String("done-location")
	if len(strings.TrimSpace(doneLocation)) == 0 {
		doneLocation = "nil"
	}
	log.Info(fmt.Sprintf("done-location : %v", doneLocation))

	doneLocation = c.String("dealer.schedule.http.done-location")
	if len(strings.TrimSpace(doneLocation)) == 0 {
		doneLocation = "nil"
	}
	log.Info(fmt.Sprintf("done-location yaml : %v", doneLocation))
}
