# ylog

## Features

1. 封装zap日志打印作为file和console的打印系统
2. 只需实现`Handler`接口，可以自由选择将日志打印到不同的系统（本地、远程）
3. 可选系统默认实例或自己组装不同的`Handler`形成不同的`logger`实例，进行不同需求的输出

#### 名词说明：

###### Handler

实现了以下方法的接口：

```go
// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services.
type Handler interface {
	Log(context.Context, Level, string, ...interface{})
	Close() error
}
```

###### logger

组合了多个`Handler`的一个日志打印实例

```go
// 结构体
type YLogger struct {
    h *Handlers
}

// 系统默认初始化一个fileHandler，注册成logger
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

// 系统提供的一系列方法
func DebugC(ctx context.Context, format string, args ...interface{})

func InfoC(ctx context.Context, format string, args ...interface{})

func WarnC(ctx context.Context, format string, args ...interface{})

func ErrorC(ctx context.Context, format string, args ...interface{})

func PanicC(ctx context.Context, format string, args ...interface{})

func FatalC(ctx context.Context, format string, args ...interface{})
```

###### Config配置

```go
type Config struct {
	FileConfig *FileConfig
}

type FileConfig struct {
    LogFilePath string `yaml:"path"`         // 日志路径
    MaxSize     int    `yaml:"max_size"`     // 单个日志最大的文件大小. 单位: MB
    MaxBackups  int    `yaml:"max_backups"`  // 日志文件最多保存多少个备份
    MaxAge      int    `yaml:"max_age"`      // 文件最多保存多少天
    Console     bool   `yaml:"console"`      // 是否命令行输出，开发环境可以使用
    LevelString string `yaml:"level_string"` // 输出的日志级别, 值：debug,info,warn,error,panic,fatal
}
```

###### 打印日志的含义

```json
{"level":"info","time":"2019-05-29T11:02:31+08:00","msg":"this is info msg, hello world","path":"ll.go:11(ll.LogTest)
  ","tid":"xxx""}
// level：日志级别
// time：时间
// msg：实际的消息
// path：打点位置：文件名:行号(包名.函数名)
// tid：trace_id
```



## 使用方法:

### 使用系统内部函数（日常默认使用这个就行）

事例：

```go
func TestDebug(t *testing.T) {
    Init(nil)
    ctx := context.WithValue(context.Background(), TraceIDKey, "trace_xxxx")
    InfoC(ctx, "this is info msg, hello %s", "world")
    ErrorC(ctx, "this is info msg, hello %s", "world")
    DebugC(ctx, "this is info msg, hello %s", "world")
    WarnC(ctx, "this is info msg, hello %s", "world")
}
```

打印如下：

```json
{"level":"info","time":"2021-07-22T12:07:00+08:00","path":"log_test.go:11(ylog.TestDebug)","tid":null,"msg":"this is info msg, hello world"}
{"level":"error","time":"2021-07-22T12:07:00+08:00","path":"log_test.go:12(ylog.TestDebug)","tid":null,"msg":"this is info msg, hello world"}
{"level":"debug","time":"2021-07-22T12:07:00+08:00","path":"log_test.go:13(ylog.TestDebug)","tid":null,"msg":"this is info msg, hello world"}
{"level":"warn","time":"2021-07-22T12:07:00+08:00","path":"log_test.go:14(ylog.TestDebug)","tid":null,"msg":"this is info msg, hello world"}

```

