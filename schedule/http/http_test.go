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
	var source = `"Content-Type=text/xml; charset=utf-8"`
	key, value, err := analysisHeader(source)
	if err != nil {
		t.Logf("errors : %s", err)
		t.Fail()
		return
	}
	t.Logf("analysis success : key[%s], value[%s]", key, value)
}
