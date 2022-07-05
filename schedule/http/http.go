package http

import (
	"context"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"sync"
	"time"
)

type Http struct {
	client *http.Client
	before []func(request *HttpRequest) error
	after  []func(response *HttpResponse) error
	mutex  *sync.Mutex
}

func (httpClient *Http) Before(beforeFunc func(request *HttpRequest) error) {
	httpClient.mutex.Lock()
	defer httpClient.mutex.Unlock()
	if httpClient.before == nil {
		httpClient.before = make([]func(request *HttpRequest) error, 16)
	}
	httpClient.before = append(httpClient.before, beforeFunc)
}

func (httpClient *Http) Handle(request *HttpRequest) error {
	if httpClient.before != nil {
		for _, beforeFunc := range httpClient.before {
			if err := beforeFunc(request); err != nil {
				log.Error(fmt.Sprintf("[dealer_cli.schedule.http.Http.Handle] do before function, request[%v], errors : %s", *request, err))
				return err
			}
		}
	}

	if httpClient.client == nil {
		return errors.New("[dealer_cli.schedule.http.Http.Handle] client not initialized ")
	}
	resp, err := httpClient.client.Do(request.Request)
	defer resp.Body.Close()
	if err != nil {
		log.Error(fmt.Sprintf("[dealer_cli.schedule.http.Http.Handle] do request[%v], errors : %s", *request, err))
		return err
	}
	response := &HttpResponse{
		resp,
		request.Source,
	}
	if httpClient.after != nil {
		for _, afterFunc := range httpClient.after {
			if err = afterFunc(response); err != nil {
				log.Error(fmt.Sprintf("[dealer_cli.schedule.http.Http.Handle] do after function, request[%v], errors : %s", *request, err))
				return err
			}
		}
	}

	return nil
}

func (httpClient *Http) After(afterFunc func(response *HttpResponse) error) {
	httpClient.mutex.Lock()
	defer httpClient.mutex.Unlock()
	if httpClient.after == nil {
		httpClient.after = make([]func(response *HttpResponse) error, 16)
	}
	httpClient.after = append(httpClient.after, afterFunc)
}

func New() *Http {
	var httpClient = &Http{
		mutex: &sync.Mutex{},
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: createTransportDialContext(&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}),
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
	}
	return httpClient
}
func createTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		log.Debug(fmt.Sprintf("dial context : %#v \n", ctx))
		return dialer.DialContext(ctx, network, addr)
	}
}
