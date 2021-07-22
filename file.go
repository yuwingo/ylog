package ylog

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type FileHandler struct {
	logger *zap.SugaredLogger
	level  Level
}

type FileConfig struct {
	LogFilePath string `mapstructure:"path"`         // 日志路径
	MaxSize     int    `mapstructure:"max_size"`     // 单个日志最大的文件大小. 单位: MB
	MaxBackups  int    `mapstructure:"max_backups"`  // 日志文件最多保存多少个备份
	MaxAge      int    `mapstructure:"max_age"`      // 文件最多保存多少天
	Console     bool   `mapstructure:"console"`      // 是否命令行输出，开发环境可以使用
	LevelString string `mapstructure:"level_string"` // 输出的日志级别, 值：debug,info,warn,error,panic,fatal
}

func NewFileHandler(c *FileConfig) *FileHandler {
	xLogTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	hookInfo := lumberjack.Logger{
		Filename:   c.LogFilePath + "info.log",
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
	}

	hookError := lumberjack.Logger{
		Filename:   c.LogFilePath + "error.log",
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
	}

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	jsonErr := zapcore.AddSync(&hookError)
	jsonInfo := zapcore.AddSync(&hookInfo)

	// Optimize the xLog output for machine consumption and the console output
	// for human operators.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = ""
	encoderConfig.TimeKey = "time"
	//encoderConfig.CallerKey = "path" // 原定的path字段含义太多，建议还是分开，然后log调用的地方就叫caller
	encoderConfig.EncodeTime = xLogTimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	xLogEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var allCore []zapcore.Core
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the cores together.
	if c.LogFilePath != "" {
		errCore := zapcore.NewCore(xLogEncoder, jsonErr, errPriority)
		infoCore := zapcore.NewCore(xLogEncoder, jsonInfo, infoPriority)
		allCore = append(allCore, errCore, infoCore)
	}

	if c.Console {
		consoleDebugging := zapcore.Lock(os.Stdout)
		allCore = append(allCore, zapcore.NewCore(xLogEncoder, consoleDebugging, infoPriority))
	}

	core := zapcore.NewTee(allCore...)

	var opts []zap.Option

	logger := zap.New(core).WithOptions(opts...).Sugar()
	defer logger.Sync()

	return &FileHandler{logger: logger, level: LevelStringToCode(c.LevelString)}
}

func log(logger *zap.SugaredLogger, l Level, keysAndValues []interface{}) {
	switch l {
	case DebugLevel:
		logger.Debugw("", keysAndValues...)
	case InfoLevel:
		logger.Infow("", keysAndValues...)
	case WarnLevel:
		logger.Warnw("", keysAndValues...)
	case ErrorLevel:
		logger.Errorw("", keysAndValues...)
	case PanicLevel:
		logger.Panicw("", keysAndValues...)
	case FatalLevel:
		logger.Fatalw("", keysAndValues...)
	}
}

func (fh *FileHandler) Log(ctx context.Context, l Level, format string, args ...interface{}) {
	if l < fh.level {
		return
	}

	logger := fh.getLogger()

	traceID := ctx.Value(TraceIDHeaderKey)
	logger = logger.With(TraceIDKey, traceID)

	msg := format
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	var keysAndValues []interface{}

	keysAndValues = append(keysAndValues, "msg")
	keysAndValues = append(keysAndValues, msg)

	log(logger, l, keysAndValues)
}

func (fh *FileHandler) Close() (err error) {
	return
}

func getCaller(skip int) string {
	fileName, line, funcName := "???", 0, "???"
	pc, fileName, line, ok := runtime.Caller(skip)
	if ok {
		funcName = runtime.FuncForPC(pc).Name() // main.(*MyStruct).foo
		funcName = filepath.Base(funcName)      // .foo
		//funcName = strings.TrimPrefix(funcName, ".") // foo

		fileName = filepath.Base(fileName) // /full/path/basename.go => basename.go
	}

	ca := fileName + ":" + strconv.Itoa(line) + "(" + funcName + ")"
	return ca
}

func (fh *FileHandler) getLogger() (logger *zap.SugaredLogger) {
	logger = fh.logger.With("path", getCaller(5))
	return
}
