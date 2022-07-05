package http

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"path/filepath"
	"regexp"
)

var location, _ = os.Getwd()
var headerFlag = &cli.StringSliceFlag{
	Name:     "header",
	Usage:    "specify the headers of http message",
	FilePath: filepath.Join(location, "./schedule.http.header.dealer"),
}

var headerYamlFlag = altsrc.NewStringSliceFlag(
	&cli.StringSliceFlag{
		Name:  "dealer.schedule.http.header",
		Usage: "specify the headers of http message. same as the flag 'header' ",
	},
)

var methodFlag = &cli.StringFlag{
	Name:  "method",
	Usage: "indicates the http method",
	Value: "GET",
}
var methodYamlFlag = altsrc.NewStringFlag(
	&cli.StringFlag{
		Name:  "dealer.schedule.http.method",
		Usage: "indicates the http method, same as the flag ‘method' ",
	},
)

var bodyFlag = &cli.StringFlag{
	Name:     "body",
	Usage:    "the body of http message. ",
	FilePath: filepath.Join(location, "./schedule.http.body.dealer"),
}

var bodyYamlFlag = altsrc.NewStringFlag(
	&cli.StringFlag{
		Name:  "dealer.schedule.http.body",
		Usage: "the body of http message. same as the flag 'body'",
	},
)

var urlFlag = &cli.StringFlag{
	Name:  "url",
	Usage: "the target url that the http message should send to",
}
var urlYamlFlag = altsrc.NewStringFlag(
	&cli.StringFlag{
		Name:  "dealer.schedule.http.url",
		Usage: "the target url that the http message should send to. same as the flag 'url' ",
	},
)

var headerFilePathRegexExpr = `"([\s\S]+?)=([\s\S]+?)"`

func analysisHeaderFilePath(value string) (map[string][]string, error) {
	headerFilePathRegex, err := regexp.Compile(headerFilePathRegexExpr)
	if err != nil {
		return nil, err
	}

	submatches := headerFilePathRegex.FindAllStringSubmatch(value, -1)
	re := make(map[string][]string, len(submatches))
	for _, submatch := range submatches {
		if len(submatch) != 3 {
			return nil, errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", value, headerFilePathRegexExpr))
		}
		if re[submatch[1]] == nil {
			re[submatch[1]] = make([]string, 0, 8)
		}
		re[submatch[1]] = append(re[submatch[1]], submatch[2])
	}
	return re, nil
}

var headerRegexExpr = `([\s\S]+?)=([\s\S]+)`

func analysisHeader(headers []string) (map[string][]string, error) {
	headerFilePathRegex, err := regexp.Compile(headerRegexExpr)
	if err != nil {
		return nil, err
	}
	re := make(map[string][]string)
	for _, header := range headers {
		submatches := headerFilePathRegex.FindAllStringSubmatch(header, -1)
		if len(submatches) != 1 {
			return nil, errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", header, headerFilePathRegexExpr))
		}
		if len(submatches[0]) != 3 {
			return nil, errors.New(fmt.Sprintf("[%s] splitted failed by regex[%s]", header, headerFilePathRegexExpr))
		}
		key, value := submatches[0][1], submatches[0][2]
		if re[key] == nil {
			re[key] = make([]string, 0, 8)
		}
		re[key] = append(re[key], value)
	}
	return re, nil
}
