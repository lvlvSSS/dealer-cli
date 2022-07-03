package schedule

import (
	"github.com/urfave/cli/v2"
	"sync"
)

func NewHttpRemote(c *cli.Context) (Remote, error) {
	return &Http{}, nil
}

type Http struct {
}

func (http *Http) Handle(wg *sync.WaitGroup) error {
	return nil
}

type HttpBuilder struct {
	method  string
	headers map[string]string
	body    []byte
}

func (builder HttpBuilder) Method(method string) HttpBuilder {
	builder.method = method
	return builder
}
