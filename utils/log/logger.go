package log

import (
	"dealer-cli/docs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

var DefaultLogName = "./Logs/dealer-cli.log"
var DefaultLevel zapcore.Level = zapcore.InfoLevel

var logger *zap.Logger

func Init(logPath string, level zapcore.Level, options ...zap.Option) {
	logger = NewLogger(logPath, level, options...)
}

func generalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

// NewLogger return a new logger that specified the log path, level, and some customized options.
func NewLogger(logPath string, level zapcore.Level, options ...zap.Option) *zap.Logger {
	if len(strings.TrimSpace(logPath)) == 0 {
		logPath = DefaultLogName
	}
	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		level = DefaultLevel
	}
	lumberjackLogger := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
		LocalTime:  true,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "log",
		CallerKey:        "linenum",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding + "----------------------------------------------------------------" + zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder, // use the lower case to output level
		EncodeTime:       generalTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " |",
	}
	// set the log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                                    // encoding configuration
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjackLogger)), // output to stdout and lumberjackLogger
		atomicLevel, // log level
	)

	// redefine the error output, to stderr and lumberjackLogger
	errorOutput := zap.ErrorOutput(zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stderr), zapcore.AddSync(&lumberjackLogger)))

	zapOptions := make([]zap.Option, 0, len(options)+1)
	for _, op := range options {
		zapOptions = append(zapOptions, op)
	}
	zapOptions = append(zapOptions, errorOutput)
	logger = zap.New(core, zapOptions...)

	return logger.Named(docs.APP_NAME)
}

type LogError struct {
	origin error
	Msg    string
}

func (err LogError) Error() string {
	if len(strings.TrimSpace(err.Msg)) == 0 && err.origin != nil {
		return err.origin.Error()
	}
	return err.Msg
}

func Debug(msg string, fields ...zap.Field) error {
	if logger == nil {
		return LogError{
			origin: nil,
			Msg:    "logger needs to be initialized first ...",
		}
	}
	logger.Debug(msg, fields...)
	return nil
}

func Info(msg string, fields ...zap.Field) error {
	if logger == nil {
		return LogError{
			origin: nil,
			Msg:    "logger needs to be initialized first ...",
		}
	}
	logger.Info(msg, fields...)
	return nil
}

func Warn(msg string, fields ...zap.Field) error {
	if logger == nil {
		return LogError{
			origin: nil,
			Msg:    "logger needs to be initialized first ...",
		}
	}
	logger.Warn(msg, fields...)
	return nil
}

func Error(msg string, fields ...zap.Field) error {
	if logger == nil {
		return LogError{
			origin: nil,
			Msg:    "logger needs to be initialized first ...",
		}
	}
	logger.Error(msg, fields...)
	return nil
}
