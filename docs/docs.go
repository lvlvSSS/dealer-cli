package docs

import (
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

const APP_NAME = "dealer-cli"

const APP_LOAD_YAML = "load-yaml"

var root, _ = os.Getwd()
var LoadFlag = &cli.StringFlag{
	Name:  APP_LOAD_YAML,
	Value: filepath.Join(root, "configs.yaml"),
}
