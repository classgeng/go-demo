package logger

import (
	"log"
	"os"
)

var std Logger = log.New(os.Stdout, "tce sdk", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

// Logger 日志日志输出接口
type Logger interface {
	Printf(format string, args ...interface{})
}

// SetLogger 设置日志接口
func SetLogger(log Logger) {
	std = log
}

// Printf ...
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
