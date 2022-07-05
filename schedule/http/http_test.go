package http

import (
	"testing"
)

func TestAnalysisHeaderFilePath(t *testing.T) {
	var source = `"Content-Type=text/xml; charset=utf-8"
"Accept-Encoding=gzip, deflate"`
	path, err := analysisHeaderFilePath(source)
	if err != nil {
		t.Logf("errors : %s", err)
		t.Fail()
		return
	}
	t.Logf("analysis success : %v", path)
}

func TestAnalysisHeader(t *testing.T) {
	headers := make([]string, 0, 8)
	headers = append(headers, `Content-Type=text/xml; charset=utf-8`)
	headers = append(headers, `Accept-Encoding=gzip, deflate`)
	headers = append(headers, `Content-Type=application/xml; charset=utf-8`)
	re, err := analysisHeader(headers)
	if err != nil {
		t.Logf("errors : %s", err)
		t.Fail()
		return
	}
	if len(re["Content-Type"]) != 2 {
		t.Fail()
		return
	}
	t.Logf("analysis success : result [%v]", re)
}
