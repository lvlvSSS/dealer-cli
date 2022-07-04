package http

import (
	"sync"
)

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
