package utility

import (
	"log"
	"io"
	"sync"
	"os"
	"runtime"
	"fmt"
	"path/filepath"
	"reflect"
)

type LogLevel uint16
const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type LOG struct {
	logger *log.Logger
	loglevel LogLevel
	mutex sync.Mutex
}
var logging = &LOG{loglevel : INFO}

// Default logger, print the log to screen.
func GetLogger() *LOG {
	if logging.logger == nil {
		logging.mutex.Lock()
		defer logging.mutex.Unlock()

		logging.logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return logging
}

func log_sprintf(format string, a ...interface{}) {
	result := ""
	for _, item := range a {
		temp := ""
		switch reflect.ValueOf(item).Kind() {
		case reflect.Bool:
			temp = fmt.Sprintf("%t", reflect.ValueOf(item).Bool())
		case reflect.Int:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Int())
		case reflect.Int8:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Int())
		case reflect.Int16:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Int())
		case reflect.Int32:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Int())
		case reflect.Int64:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Int())
		case reflect.Uint:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Uint())
		case reflect.Uint8:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Uint())
		case reflect.Uint16:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Uint())
		case reflect.Uint32:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Uint())
		case reflect.Uint64:
			temp = fmt.Sprintf("%d", reflect.ValueOf(item).Uint())
		case reflect.Uintptr:
			temp = fmt.Sprintf("%p", reflect.ValueOf(item).Pointer())
		case reflect.Float32:
			temp = fmt.Sprintf("%f", reflect.ValueOf(item).Float())
		case reflect.Float64:
			temp = fmt.Sprintf("%f", reflect.ValueOf(item).Float())
		case reflect.Complex64:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Complex())
		case reflect.Complex128:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Complex())
		case reflect.Array:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Chan:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Func:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Interface:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Map:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Ptr:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.Slice:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		case reflect.String:
			temp = fmt.Sprintf("%s", reflect.ValueOf(item).String())
		case reflect.Struct:
			temp = fmt.Sprintf("%s", reflect.ValueOf(item).Interface())
		default:
			temp = fmt.Sprintf("%v", reflect.ValueOf(item).Interface())
		}

		result = result + temp
	}

}

// To check if the log level is reach the setting
func reachable(loggerLevel LogLevel, funcLevel LogLevel) bool {
	if funcLevel  >= loggerLevel {
		return true
	}

	return false
}

func getCallerInfo() (fn string, line int) {
	_, fn, line, _ = runtime.Caller(2)
	fn = filepath.Base(fn)

	return fn,line
}

func (l *LOG) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

func (l *LOG) SetMinorLogLevel(level LogLevel) {
	l.loglevel = level
}

func (l *LOG) DEBUG(arg string) {
	fn, line := getCallerInfo()
	tmp := fmt.Sprintf("[%s:%d]  %s", fn, line, arg)

	// If log level is not reach the setting, do not print it.
	if ! reachable(l.loglevel, DEBUG) {
		return
	}

	l.logger.SetPrefix("[DEBUG]")
	l.logger.Println(tmp)
}

func (l *LOG) INFO(arg string) {
	fn, line := getCallerInfo()
	tmp := fmt.Sprintf("[%s:%d]  %s", fn, line, arg)

	// If log level is not reach the setting, do not print it.
	if ! reachable(l.loglevel, INFO) {
		return
	}

	l.logger.SetPrefix("[INFO]")
	l.logger.Println(tmp)
}

func (l *LOG) WARN(arg string) {
	fn, line := getCallerInfo()
	tmp := fmt.Sprintf("[%s:%d]  %s", fn, line, arg)

	// If log level is not reach the setting, do not print it.
	if ! reachable(l.loglevel, WARN) {
		return
	}

	l.logger.SetPrefix("[WARN]")
	l.logger.Println(tmp)
}

func (l *LOG) ERROR(arg string) {
	fn, line := getCallerInfo()
	tmp := fmt.Sprintf("[%s:%d]  %s", fn, line, arg)

	// If log level is not reach the setting, do not print it.
	if ! reachable(l.loglevel, ERROR) {
		return
	}

	l.logger.SetPrefix("[ERROR]")
	l.logger.Println(tmp)
}

func (l *LOG) TRACE(arg string) {
	fn, line := getCallerInfo()
	tmp := fmt.Sprintf("[%s:%d]  %s", fn, line, arg)

	// If log level is not reach the setting, do not print it.
	if ! reachable(l.loglevel, TRACE) {
		return
	}

	l.logger.SetPrefix("[TRACE]")
	l.logger.Println(tmp)
}