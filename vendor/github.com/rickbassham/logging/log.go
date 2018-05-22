package logging

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarning = "WARN"
	LogLevelError   = "ERROR"
)

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger

	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type logger struct {
	f     Formatter
	o     io.Writer
	level string
}

type errorLocation struct {
	Package  string `json:"package"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type logEntry struct {
	l *logger

	Timestamp     time.Time              `json:"timestamp,omitempty"`
	Level         string                 `json:"level"`
	Package       string                 `json:"package"`
	Function      string                 `json:"function"`
	File          string                 `json:"file"`
	Line          int                    `json:"line"`
	Fields        map[string]interface{} `json:"fields,omitempty"`
	ErrorStr      string                 `json:"error,omitempty"`
	ErrorData     error                  `json:"errorData,omitempty"`
	ErrorLocation *errorLocation         `json:"errorLocation,omitempty"`
	Msg           string                 `json:"message"`
}

type Formatter interface {
	Format(*logEntry) (string, error)
}

type JSONFormatter struct{}

func NewLogger(output io.Writer, f Formatter, level string) Logger {
	if level != LogLevelDebug && level != LogLevelInfo && level != LogLevelWarning && level != LogLevelError {
		level = LogLevelInfo
	}

	if output == nil {
		output = ioutil.Discard
	}

	if f == nil {
		f = JSONFormatter{}
	}

	return &logger{
		o:     output,
		f:     f,
		level: level,
	}
}

func (l *logger) WithField(key string, value interface{}) Logger {
	le := logEntry{
		l: l,
		Fields: map[string]interface{}{
			key: value,
		},
	}

	return &le
}

func (l *logger) WithError(err error) Logger {
	le := logEntry{
		l:         l,
		Fields:    map[string]interface{}{},
		ErrorData: err,
	}

	if err != nil {
		le.ErrorStr = err.Error()
	}

	le.setErrorLocation()

	return &le
}

func (l *logger) Debug(msg string) {
	if l.level == LogLevelError || l.level == LogLevelWarning || l.level == LogLevelInfo {
		return
	}

	le := &logEntry{
		l:      l,
		Fields: map[string]interface{}{},
	}

	le.setData(msg, LogLevelDebug)
	le.write()
}

func (l *logger) Info(msg string) {
	if l.level == LogLevelError || l.level == LogLevelWarning {
		return
	}

	le := &logEntry{
		l:      l,
		Fields: map[string]interface{}{},
	}

	le.setData(msg, LogLevelInfo)
	le.write()
}

func (l *logger) Warn(msg string) {
	if l.level == LogLevelError {
		return
	}

	le := &logEntry{
		l:      l,
		Fields: map[string]interface{}{},
	}

	le.setData(msg, LogLevelWarning)
	le.write()
}

func (l *logger) Error(msg string) {
	le := &logEntry{
		l:      l,
		Fields: map[string]interface{}{},
	}

	le.setData(msg, LogLevelError)
	le.write()
}

func (l *logEntry) newEntry() *logEntry {
	e := logEntry{
		l:         l.l,
		Fields:    map[string]interface{}{},
		ErrorData: l.ErrorData,
		ErrorStr:  l.ErrorStr,
	}

	for k, v := range l.Fields {
		e.Fields[k] = v
	}

	return &e
}

func (l *logEntry) WithField(key string, value interface{}) Logger {
	e := l.newEntry()
	e.Fields[key] = value
	return e
}

func (l *logEntry) WithError(err error) Logger {
	e := l.newEntry()

	e.ErrorData = err
	if err != nil {
		e.ErrorStr = err.Error()
	}

	l.setErrorLocation()

	return e
}

func (l *logEntry) setErrorLocation() {
	loc := retrieveCallInfo()
	l.ErrorLocation = &errorLocation{
		File:     loc.fileName,
		Function: loc.funcName,
		Line:     loc.line,
		Package:  loc.packageName,
	}
}

func (l *logEntry) setData(msg, level string) {
	caller := retrieveCallInfo()

	l.File = caller.fileName
	l.Line = caller.line
	l.Function = caller.funcName
	l.Package = caller.packageName
	l.Msg = msg
	l.Level = level
	l.Timestamp = time.Now()
}

func (l *logEntry) Debug(msg string) {
	if l.l.level == LogLevelError || l.l.level == LogLevelWarning || l.l.level == LogLevelInfo {
		return
	}

	l.setData(msg, LogLevelDebug)
	l.write()
}

func (l *logEntry) Info(msg string) {
	if l.l.level == LogLevelError || l.l.level == LogLevelWarning {
		return
	}

	l.setData(msg, LogLevelInfo)
	l.write()
}

func (l *logEntry) Warn(msg string) {
	if l.l.level == LogLevelError {
		return
	}

	l.setData(msg, LogLevelWarning)
	l.write()
}

func (l *logEntry) Error(msg string) {
	l.setData(msg, LogLevelError)
	l.write()
}

func (l *logEntry) write() {
	output, err := l.l.f.Format(l)
	if err != nil {
		l.l.o.Write([]byte("error marshalling logEntry"))
		return
	}

	l.l.o.Write([]byte(output))
	l.l.o.Write([]byte("\n"))
}

func (JSONFormatter) Format(l *logEntry) (string, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(3)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
