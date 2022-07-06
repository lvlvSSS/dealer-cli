package http

import (
	"dealer-cli/utils/converter"
	file_util "dealer-cli/utils/files"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
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
	url        string
	method     string
	headers    map[string][]string
	bodyFormat string

	doneLocation  string // file is sent, then move the file to doneLocation directory.
	fileFormat    string
	fileSourceDir string

	requests chan *HttpRequest
	mutex    *sync.Mutex
}

var doneLocationFlag = &cli.StringFlag{
	Name:  "done-location",
	Usage: "file is sent, then move the file to the directory.",
	Value: filepath.Join(location, defaultLocation),
}

var doneLocationYamlFlag = altsrc.NewStringFlag(
	&cli.StringFlag{
		Name:  "dealer.schedule.http.done-location",
		Usage: "file is sent, then move the file to the directory. same as the flag 'done-location' ",
	},
)

var defaultLocation = "./done"
var fileFormatRegexPattern = `#\{(\S+?)\}`

func (producer *FileRequestProducer) Init(c *cli.Context) error {
	// set url
	if producer.url = c.String("url"); len(strings.TrimSpace(producer.url)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] url[%s]", producer.url))
	} else if producer.url = c.String("dealer.schedule.http.url"); len(strings.TrimSpace(producer.url)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] dealer.schedule.http.url[%s]", producer.url))
	} else {
		log.Error("[dealer_cli.schedule.http.FileRequestProducer.Init] url is empty")
		return errors.New("[dealer_cli.schedule.http.FileRequestProducer.Init] url is empty")
	}
	// set method
	if c.IsSet("method") {
		producer.method = strings.ToUpper(c.String("method"))
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] method[%s]", producer.method))
	} else if producer.method = strings.ToUpper(c.String("dealer.schedule.http.method")); len(strings.TrimSpace(producer.method)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] dealer.schedule.http.method[%s]", producer.method))
	} else {
		producer.method = strings.ToUpper(c.String("method"))
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] method forced to [%s]", producer.method))
	}
	// set headers
	if headerSource := c.StringSlice("dealer.schedule.http.header"); len(headerSource) > 0 {
		headers, err := analysisHeader(headerSource)
		if err == nil {
			producer.headers = headers
		}
	}
	if headerSource := c.StringSlice("header"); len(headerSource) > 0 {
		// if length is 1, the header may be from the file[schedule.http.header.dealer]
		if len(headerSource) == 1 {
			headers, err := analysisHeaderFilePath(headerSource[0])
			if err == nil {
				producer.headers = headers
			}
		} else {
			headers, err := analysisHeader(headerSource)
			if err == nil {
				producer.headers = headers
			}
		}
	}
	// set bodyFormat
	if producer.bodyFormat = c.String("body"); len(strings.TrimSpace(producer.bodyFormat)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] bodyFormat[%s]", producer.bodyFormat))
	} else if producer.bodyFormat = c.String("dealer.schedule.http.body"); len(strings.TrimSpace(producer.bodyFormat)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] dealer.schedule.http.body[%s]", producer.bodyFormat))
	} else {
		log.Error("[dealer_cli.schedule.http.FileRequestProducer.Init] bodyFormat is empty")
		return errors.New("[dealer_cli.schedule.http.FileRequestProducer.Init] bodyFormat is empty")
	}
	// check fileSourceDir and fileFormat
	if err := producer.compileBody(fileFormatRegexPattern); err != nil {
		return err
	}
	// set doneLocation
	if c.IsSet("done-location") {
		producer.doneLocation = c.String("done-location")
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] done-location[%s]", producer.doneLocation))
	} else if producer.doneLocation = c.String("dealer.schedule.http.done-location"); len(strings.TrimSpace(producer.doneLocation)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] dealer.schedule.http.done-location[%s]", producer.doneLocation))
	} else {
		producer.doneLocation = c.String("done-location")
		log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] default done-location[%s]", producer.doneLocation))
	}
	doneLocationAbs, err := filepath.Abs(producer.doneLocation)
	if err != nil {
		return err
	}
	producer.doneLocation = doneLocationAbs
	locStat, err := os.Stat(producer.doneLocation)
	if err != nil && !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] init failed : doneLocation[%s] is not directory", producer.doneLocation))
	} else if err != nil && os.IsNotExist(err) {
		os.MkdirAll(producer.doneLocation, 0)
	} else if !locStat.IsDir() {
		return errors.New(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Init] init failed : doneLocation[%s] is not directory", producer.doneLocation))
	}

	producer.mutex = &sync.Mutex{}

	return nil
}

func (producer *FileRequestProducer) compileBody(regexExpr string) error {
	fileFormatRegex, err := regexp.Compile(regexExpr)
	if err != nil {
		return err
	}
	submatches := fileFormatRegex.FindAllStringSubmatch(producer.bodyFormat, 1)
	if len(submatches) != 1 || len(submatches[0]) != 2 {
		return errors.New(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.compileBody] compile body[%s] by regex[%s] failed ...", producer.bodyFormat, regexExpr))
	}
	fullPath, err := filepath.Abs(submatches[0][1])
	if err != nil {
		return err
	}
	producer.fileFormat = fullPath
	producer.fileSourceDir = filepath.Clean(strings.TrimRight(fullPath, filepath.Base(producer.fileFormat)))
	fileSourceStat, err := os.Stat(producer.fileSourceDir)
	if err != nil || !fileSourceStat.IsDir() {
		return errors.New(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.compileBody] compile failed : fileSourceDir[%s] is not directory", producer.fileSourceDir))
	}
	producer.bodyFormat = fileFormatRegex.ReplaceAllString(producer.bodyFormat, "%s")

	return nil
}

func (producer *FileRequestProducer) Stream() <-chan *HttpRequest {
	producer.mutex.Lock()
	defer producer.mutex.Unlock()
	if producer.requests != nil {
		return producer.requests
	}

	requests := make(chan *HttpRequest, runtime.NumCPU()*2+1)
	producer.requests = requests

	files, _ := file_util.GetAllSubFiles(producer.fileSourceDir)
	if len(files) <= (runtime.NumCPU()*2 + 1) {
		for _, file := range files {
			if req := producer.isValid(file); req != nil {
				log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Stream] file[%s] stack in ...", file))
				requests <- req
			}
		}
		return requests
	}
	count := 0
	index := 0
	for {
		if count >= (runtime.NumCPU()*2+1) || index >= len(files) {
			break
		}
		if req := producer.isValid(files[index]); req != nil {
			log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Stream] file[%s] stack in ...", files[index]))
			requests <- req
			count++
		}
		index++
	}
	if index < len(files) {
		go func() {
			for _, file := range files[index:] {
				if req := producer.isValid(file); req != nil {
					log.Debug(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.Stream] file[%s] stack in ...", file))
					requests <- req
				}
			}
		}()
	}
	return requests
}

func (producer *FileRequestProducer) isValid(file string) *HttpRequest {
	isMatch, err := filepath.Match(producer.fileFormat, file)
	if err != nil {
		log.Error(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.isValid] match file[%s] errors : %v", file, err))
		return nil
	}
	if !isMatch {
		log.Error(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.isValid] file[%s] doesn't match format[%s] ", file, producer.fileFormat))
		return nil
	}
	source, err := os.ReadFile(file)
	if err != nil {
		log.Error(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.isValid] read file[%s] errors : %v", file, err))
		return nil
	}

	req, err := http.NewRequest(producer.method, producer.url, strings.NewReader(fmt.Sprintf(producer.bodyFormat, converter.BytesToString(source))))
	if err != nil {
		log.Error(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.isValid] create request failed[%s] errors : %v", file, err))
		return nil
	}
	// set headers
	if producer.headers != nil {
		for key, values := range producer.headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	return &HttpRequest{
		Request: req,
		Source:  file,
	}
}

func (producer *FileRequestProducer) After(response *HttpResponse) error {
	all, err := ioutil.ReadAll(response.Response.Body)
	if err == nil {
		log.Info(
			fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.After] request from file[%s], response is [%s]",
				response.Source,
				converter.BytesToString(all)))
	}

	if response.StatusCode != 200 {
		log.Warn(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.After] send file[%s] failed.", response.Source))
		return nil
	}
	if err := os.Rename(response.Source, filepath.Join(producer.doneLocation, filepath.Base(response.Source))); err != nil {
		log.Warn(fmt.Sprintf("[dealer_cli.schedule.http.FileRequestProducer.After] rename from %s to %s fails, error[%v]",
			response.Source, filepath.Join(producer.doneLocation, filepath.Base(response.Source)), err))
	}
	return nil
}
