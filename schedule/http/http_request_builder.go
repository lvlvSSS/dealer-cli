package http

import (
	"github.com/urfave/cli/v2"
	"net/http"
)

type HttpRequest struct {
	*http.Request
	Source string
}
type HttpResponse struct {
	*http.Response
	Source string
}

type HttpRequestBuilder interface {
	Stream() <-chan *HttpRequest // async method, return a continuous channel
}

type FileRequestProducer struct {
	method     string
	headers    map[string]string
	location   string
	fileFormat string
	doneOutput string // file is sent, then move the file to doneOutput directory.
}

func (producer *FileRequestProducer) Method(method string) *FileRequestProducer {
	producer.method = method
	return producer
}

func (producer *FileRequestProducer) Init(c *cli.Context) error {
	return nil
}

func (producer *FileRequestProducer) Stream() <-chan *HttpRequest {
	requests := make(chan *HttpRequest, 32)
	return requests
}

func (producer *FileRequestProducer) After(response *HttpResponse) error {
	return nil
}
