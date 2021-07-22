package ylog

import (
	"context"
	"testing"
)

func TestDebug(t *testing.T) {
	Init(nil)
	ctx := context.WithValue(context.Background(), TraceIDKey, "trace_xxxx")
	InfoC(ctx, "this is info msg, hello %s", "world")
	ErrorC(ctx, "this is info msg, hello %s", "world")
	DebugC(ctx, "this is info msg, hello %s", "world")
	WarnC(ctx, "this is info msg, hello %s", "world")
}

func TestManyLogger(t *testing.T) {
	loggerA := NewLogger(nil)
	ctx := context.WithValue(context.Background(), TraceIDKey, "trace_xxxx")
	loggerA.InfoC(ctx, "this is (A) info msg, hello %s", "world")
	loggerA.DebugC(ctx, "this is (A) debug msg, hello %s", "world")

	loggerB := NewLogger(&Config{FileConfig: &FileConfig{LogFilePath: "./testpath/"}})
	loggerB.DebugC(ctx, "this is (B) debug msg, hello %s", "world")
}
