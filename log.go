package ylog

import (
	"context"
)

type YLogger struct {
	h *Handlers
}

type Config struct {
	FileConfig *FileConfig `mapstructure:"fileConfig"`
}

var (
	yLogger YLogger
)

func Init(c *Config) {
	if c == nil {
		fileConfig := &FileConfig{
			Console: true,
		}
		c = &Config{FileConfig: fileConfig}
	}
	fileHandler := NewFileHandler(c.FileConfig)
	yLogger.h = NewHandlers(fileHandler)
}

func NewLogger(c *Config) *YLogger {
	var xLogger = &YLogger{}
	if c == nil {
		fileConfig := &FileConfig{
			Console: true,
		}
		c = &Config{FileConfig: fileConfig}
	}
	fileHandler := NewFileHandler(c.FileConfig)
	xLogger.SetHandlers(NewHandlers(fileHandler))
	return xLogger
}

func (xLogger *YLogger) SetHandlers(hs *Handlers) {
	xLogger.h = hs
}

func DebugC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, DebugLevel, format, args...)
}

func InfoC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, InfoLevel, format, args...)
}

func WarnC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, WarnLevel, format, args...)
}

func ErrorC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, ErrorLevel, format, args...)
}

func PanicC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, PanicLevel, format, args...)
}

func FatalC(ctx context.Context, format string, args ...interface{}) {
	yLogger.h.Log(ctx, FatalLevel, format, args...)
}

// DebugC with context logs a message.
func (xLogger *YLogger) DebugC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, DebugLevel, format, args...)
}

// InfoC with context logs a message.
func (xLogger *YLogger) InfoC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, InfoLevel, format, args...)
}

// WarnC with context logs a message.
func (xLogger *YLogger) WarnC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, WarnLevel, format, args...)
}

// ErrorC with context logs a message.
func (xLogger *YLogger) ErrorC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, ErrorLevel, format, args...)
}

// PanicC with context logs a message.
func (xLogger *YLogger) PanicC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, PanicLevel, format, args...)
}

// FatalC with context logs a message.
func (xLogger *YLogger) FatalC(ctx context.Context, format string, args ...interface{}) {
	xLogger.h.Log(ctx, FatalLevel, format, args...)
}
