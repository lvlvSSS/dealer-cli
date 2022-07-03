package test

import (
	"dealer-cli/file/extract"
	"dealer-cli/utils/log"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type stackTrace struct {
}

func (stackTrace) Enabled(level zapcore.Level) bool {
	return true
}

// ExampleInfo - test the logger.Info
func ExampleInfo() {
	log.Init(log.DefaultLogName, zap.DebugLevel, zap.AddStacktrace(stackTrace{}), zap.AddCaller(), zap.Development())
	log.Debug("This is a testing ...")
	// Output:
	// 2022-06-16 00:09:50.965 |debug |benchmark-cli |D:/Coding/go/projects/benchmark-cli/utils/log/logger.go:108 |This is a testing ...
	// benchmark-cli/utils/log.Debug
	// D:/Coding/go/projects/benchmark-cli/utils/log/logger.go:108
	// benchmark-cli/test.ExampleInfo
	// D:/Coding/go/projects/benchmark-cli/test/example_test.go:24
	// testing.runExample
	// D:/Coding/go/go/src/testing/run_example.go:63
	// testing.runExamples
	// D:/Coding/go/go/src/testing/example.go:44
	// testing.(*M).Run
	// D:/Coding/go/go/src/testing/testing.go:1721
	// main.main
	// _testmain.go:53
	// runtime.main
	// D:/Coding/go/go/src/runtime/proc.go:250
	// ----------------------------------------------------------------
}

type testFormat struct {
	a string `number a`
}

// ExamplePrintf - test the format output
func ExamplePrintf() {
	fmt.Printf("%q %v %#v %T\n", '我', '我', '我', '我')
	a := testFormat{"abc"}
	fmt.Printf("%+v \n", a)
	// OUTPUT:
	// '我' 25105 25105 int32
	// {a:abc}
}

func ExampleExtractor() {
	log.Init(log.DefaultLogName, log.DefaultLevel)
	target := `<\?xml version="1.0" encoding="UTF-8" standalone="yes"\?>\s*<ROOT>[\s\S]+<OPTYPE>501</OPTYPE>[\s\S]+</ROOT>`
	headlineFormat := `^\[[\S\s]+\]\s\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2},\d{3}\s\[\S+?\]\s[A-Z]+\s[\s\S]+`
	getwd, _ := os.Getwd()
	fileSourcePath := filepath.Join(getwd, "log.log")
	fileTargetFormat := `#{./ROOT/CONSIS_PRESC_MSTVW/PRESC_NO} - #{./ROOT/CONSIS_PRESC_MSTVW/PATIENT_NAME}.log`
	location := filepath.Join(getwd, "output")
	extractor, _ := extract.New(target, headlineFormat, fileSourcePath, fileTargetFormat, location, 1)
	err := extractor.Extract()
	if err != nil {
		fmt.Println(err)
	}
	extractor.Done()
	// OUTPUT:
	// ss
}
