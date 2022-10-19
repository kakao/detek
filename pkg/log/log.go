package log

/*
	@scotty.scott: Except "Result()", all function will print to "stderr".
*/
import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	mu       sync.Mutex
	logKey   = &struct{}{}
	logLevel = 3
)

func EnableDebugMode() {
	logLevel = 0
}

func SetLogLevel(level string) {
	level = strings.ToLower(level)
	v, ok := map[string]int{
		"debug":  0,
		"info":   1,
		"error":  2,
		"result": 3,
	}[level]
	if !ok {
		logLevel = 1
	} else {
		logLevel = v
	}
}
func SetContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, logKey, name)
}
func NewContext(name string) context.Context {
	ctx := context.Background()
	return SetContext(ctx, name)
}
func GetContext(ctx context.Context) string {
	v := ctx.Value(logKey)
	if name, ok := v.(string); ok {
		return name
	} else {
		return ""
	}
}

func print(status, name, msg string) {
	msg = fmt.Sprintf("%s [%s]:%s\n", status, name, msg)

	mu.Lock()
	fmt.Fprint(os.Stderr, msg)
	mu.Unlock()
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	if logLevel > 0 {
		return
	}
	status := "DEBUG"
	msg = fmt.Sprintf(msg, args...)
	name := GetContext(ctx)

	print(status, name, msg)
}

func Info(ctx context.Context, msg string, args ...interface{}) {
	if logLevel > 1 {
		return
	}
	status := "INFO"
	msg = fmt.Sprintf(msg, args...)
	name := GetContext(ctx)

	print(status, name, msg)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	if logLevel > 2 {
		return
	}
	status := "ERROR"
	msg = fmt.Sprintf(msg, args...)
	name := GetContext(ctx)

	print(status, name, msg)
}
