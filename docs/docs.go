package docs

import "github.com/urfave/cli/v2"

const APP_NAME = "dealer-cli"

const APP_LOAD_YAML = "load-yaml"

var LoadFlag = &cli.StringFlag{Name: "load-yaml"}
