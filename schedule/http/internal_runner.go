package http

import (
	schedule_internal "dealer-cli/schedule/internal"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"runtime"
	"time"
)

type runner struct {
	repeat int

	duration time.Duration
	deadline <-chan time.Time

	times    int
	terminal chan int

	exit chan interface{} // signal imply that the cron should stop and program terminate.

	client *Http
	HttpRequestBuilder
}

func (run *runner) Run() {
	stream := run.Stream()

	select {
	case req, ok := <-stream:
		if !ok {
			log.Info("dealer-cli schedule http Run - stream closed ")
			schedule_internal.DefaultCron.Stop()
			return
		}
		log.Info(fmt.Sprintf("dealer-cli schedule http Run - ready to handle message[%s] ...", req.Source))
		run.client.Handle(req)
	default:
		close(run.exit)
	}

	if run.terminal != nil {
		run.terminal <- 1
	}
}

func (run *runner) Done() {
	defer schedule_internal.DefaultCron.Stop()
	if run.deadline != nil && run.terminal != nil {
		select {
		case <-run.deadline:
			log.Info(fmt.Sprintf("dealer-cli schedule http - reach the deadline[%s]", run.duration))
			break
		case _, ok := <-run.terminal:
			if !ok {
				log.Info("dealer-cli schedule http - term closed")
				return
			}
			run.times--
			if run.times <= 0 {
				log.Info(fmt.Sprintf("dealer-cli schedule http - reach the term[%d]", run.times))
				break
			}
		case <-run.exit:
			log.Info("dealer-cli schedule http - exit ...")
			return
		}
	} else if run.deadline != nil {
		select {
		case <-run.deadline:
			log.Info(fmt.Sprintf("dealer-cli schedule http - reach the deadline[%s]", run.duration))
			break
		case <-run.exit:
			log.Info("dealer-cli schedule http - exit ...")
			return
		}
	} else if run.terminal != nil {
		select {
		case _, ok := <-run.terminal:
			if !ok {
				log.Info("dealer-cli schedule http - term closed")
				return
			}
			run.times--
			if run.times <= 0 {
				log.Info(fmt.Sprintf("dealer-cli schedule http - reach the term[%d]", run.times))
				break
			}
		case <-run.exit:
			log.Info("dealer-cli schedule http - exit ...")
			return
		}
	} else {
		select {
		case <-run.exit:
			log.Info("dealer-cli schedule http - exit ...")
			return
		}
	}
}

func checkRunner(c *cli.Context, httpRequestBuilder HttpRequestBuilder, httpClient *Http) (*runner, error) {
	// check repeat
	var repeat = 0
	if c.IsSet("repeat") {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - repeat[%d]", repeat))
		repeat = c.Int("repeat")
	} else if repeat = c.Int("dealer.schedule.repeat"); repeat > 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - dealer.schedule.repeat[%d]", repeat))
	} else {
		repeat = c.Int("repeat")
	}
	if repeat <= 0 {
		log.Warn("dealer_cli schedule http - repeat forced to 1")
		repeat = 1
	}
	// check times
	var times = 0
	if times = c.Int("times"); times > 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - times[%d]", times))
	} else if times = c.Int("dealer.schedule.times"); times > 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - dealer.schedule.times[%d]", times))
	} else {
		return nil, errors.New("dealer_cli schedule http - times is not set")
	}
	// check duration
	var duration int64 = 0
	if duration = c.Int64("duration"); duration > 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - duration[%d]", duration))
	} else if duration = c.Int64("dealer.schedule.duration"); duration > 0 {
		log.Debug(fmt.Sprintf("dealer_cli schedule http - dealer.schedule.duration[%d]", duration))
	} else {
		return nil, errors.New("dealer_cli schedule http - duration is not set")
	}

	return &runner{
		repeat:             repeat,
		duration:           time.Second * time.Duration(duration),
		deadline:           time.Tick(time.Duration(duration) * time.Second),
		times:              times,
		terminal:           make(chan int, runtime.NumCPU()*2+1),
		client:             httpClient,
		exit:               make(chan interface{}),
		HttpRequestBuilder: httpRequestBuilder,
	}, nil
}
