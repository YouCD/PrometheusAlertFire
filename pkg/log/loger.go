package log

import (
	"PrometheusAlertFire/pkg/config"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

var logger *zap.SugaredLogger
var AtomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)

func init() {
	log := NewLogger(config.Cfg.Alert.LogLevel)
	// defer log.Sync()
	logger = log.Sugar()
	logger.Sync()
}

func NewLogger(level string) *zap.Logger {
	core := newCore(level)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development(), zap.AddStacktrace(zap.ErrorLevel))
}

//func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
func newCore(level string) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		// EncodeCaller:   zapcore.FullCallerEncoder,    // 全路径编码器
		EncodeName: zapcore.FullNameEncoder,
	}

	// 设置级别
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	AtomicLevel.SetLevel(logLevel)
	return zapcore.NewCore(

		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
		AtomicLevel, // 日志级别
	)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Info(err)
	}
	return dir
}

func getFilePath() string {
	logfile := getCurrentDirectory() + "/" + getAppname() + ".log"
	return logfile
}

func getAppname() string {
	full := os.Args[0]
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		return name
	}
	return ""
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}
