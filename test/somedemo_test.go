package test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHttpClient(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: createTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	backgroundContext := context.Background()
	fmt.Printf(" background context : %#v \n", backgroundContext)
	req, err := http.NewRequestWithContext(backgroundContext, "POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func createTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, s string, s2 string) (net.Conn, error) {
		fmt.Printf("dial context : %#v \n", ctx)
		return dialer.DialContext(ctx, s, s2)
	}
}

func TestFormat(t *testing.T) {
	fmt.Printf("%q \n", '我')
	fmt.Printf("%v \n", '我')
	fmt.Printf("%#v \n", '我')
	fmt.Printf("%T \n", '我')
	a := testFormat{"abc"}
	fmt.Printf("%+v \n", a)
}

func TestRoot(t *testing.T) {
	root, _ := os.Getwd()
	fmt.Println(filepath.Join(root, "./file.extract.headline.dealer"))
	info, err := os.Stat("./somedemo_test.go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(info.Size())
	}
}
