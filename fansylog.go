package fansylog

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// 自定义一个日志

// LogLevel 日志等级
type LogLevel uint16

// LogColor 日志颜色
type LogColor string

const (
	// UNKNOWN 0
	UNKNOWN LogLevel = iota
	// DEBUG 1
	DEBUG
	// TRACE 2
	TRACE
	// INFO 3
	INFO
	// WARNING 4
	WARNING
	// ERROR 5
	ERROR
	// FATAL 6
	FATAL
)

const (
	// DEBUGCOLOR ...
	DEBUGCOLOR LogColor = "\033[46m"
	// TRACECOLOR ...
	TRACECOLOR LogColor = "\033[45m"
	// INFOCOLOR ...
	INFOCOLOR LogColor = "\033[42m"
	// WARNINGCOLOR ...
	WARNINGCOLOR LogColor = "\033[43m"
	// ERRORCOLOR ...
	ERRORCOLOR LogColor = "\033[40m"
	// FATALCOLOR ...
	FATALCOLOR LogColor = "\033[41m"
)

// Logger 定义接口
type Logger interface {
	Debug(format string, a ...interface{})
	Trace(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "TRACE":
		return TRACE, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOWN, err
	}
}

func getLevelString(lv LogLevel) (string, error) {
	switch lv {
	case DEBUG:
		return "DEBUG", nil
	case TRACE:
		return "TRACE", nil
	case INFO:
		return "INFO", nil
	case WARNING:
		return "WARNING", nil
	case ERROR:
		return "ERROR", nil
	case FATAL:
		return "FATAL", nil
	default:
		err := errors.New("无效的日志级别")
		return "UNKNOWN", err
	}
}

func getLogColor(l LogLevel) (color LogColor) {
	switch l {
	case DEBUG:
		color = DEBUGCOLOR
	case TRACE:
		color = TRACECOLOR
	case INFO:
		color = INFOCOLOR
	case WARNING:
		color = WARNINGCOLOR
	case ERROR:
		color = ERRORCOLOR
	case FATAL:
		color = FATALCOLOR
	}
	return
}

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Call() failed!")
		return
	}
	fileName = path.Base(file)
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	return
}
