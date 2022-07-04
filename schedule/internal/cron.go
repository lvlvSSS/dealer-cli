package internal

import (
	"dealer-cli/utils/log"
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

var DefaultCron = cron.New(cron.WithLogger(&cronLogger{}), cron.WithChain(cron.Recover(&cronLogger{})))

type cronLogger struct{}

func (logger *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Info(fmt.Sprintf(formatString(len(keysAndValues)), append([]interface{}{msg}, keysAndValues...)...))
}

func (logger *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Error(fmt.Sprintf(formatString(len(keysAndValues)+2), append([]interface{}{msg, "error", err}, keysAndValues...)...))
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}
