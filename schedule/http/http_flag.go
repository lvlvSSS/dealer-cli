package http

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"path/filepath"
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
	})
