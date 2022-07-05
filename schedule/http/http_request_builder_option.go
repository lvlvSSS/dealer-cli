package http

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func WithMethod(method string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.method = method
	}
}

func WithHeaders(headers map[string][]string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.headers = headers
	}
}

func WithUrl(url string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.url = url
	}
}

func WithBody(body string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.bodyFormat = body
	}
}

func WithLocation(location string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.doneLocation = location
	}
}

func WithFileFormat(fileFormat string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.fileFormat = fileFormat
	}
}

func WithFileSourceDir(fileSourceDir string) func(*FileRequestProducer) {
	return func(producer *FileRequestProducer) {
		producer.fileSourceDir = fileSourceDir
	}
}

func BuildFileRequestProducer(opts ...Option) (*FileRequestProducer, error) {
	producer := &FileRequestProducer{mutex: &sync.Mutex{}}
	for _, opt := range opts {
		opt(producer)
	}
	if len(strings.TrimSpace(producer.url)) == 0 {
		return nil, errors.New("dealer_cli schedule http - build FileRequestProducer errors[url is empty]")
	}
	if len(strings.TrimSpace(producer.fileFormat)) == 0 {
		return nil, errors.New("dealer_cli schedule http - build FileRequestProducer errors[fileFormat is empty]")
	}
	if len(strings.TrimSpace(producer.fileSourceDir)) == 0 {
		return nil, errors.New("dealer_cli schedule http - build FileRequestProducer errors[fileSourceDir is empty]")
	}
	if len(strings.TrimSpace(producer.method)) == 0 {
		producer.method = "GET"
	}
	if len(strings.TrimSpace(producer.doneLocation)) == 0 {
		producer.doneLocation = defaultLocation
	}
	doneLocationAbs, err := filepath.Abs(producer.doneLocation)
	if err != nil {
		return nil, err
	}
	producer.doneLocation = doneLocationAbs
	locStat, err := os.Stat(producer.doneLocation)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("dealer_cli schedule http -  build FileRequestProducer failed : doneLocation[%s] is not directory", producer.doneLocation))
	} else if err != nil && os.IsNotExist(err) {
		os.MkdirAll(producer.doneLocation, 0)
	} else if !locStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("dealer_cli schedule http -  build FileRequestProducer failed : doneLocation[%s] is not directory", producer.doneLocation))
	}

	return producer, nil
}

type Option func(*FileRequestProducer)
