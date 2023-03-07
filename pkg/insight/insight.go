package insight

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger func(logLevel LogLevels, id string, data any)

type LogLevels string

var (
	_logger     Logger
	_middleware []func(id string, origin string, logLevel LogLevels, fields map[string]any)
)

const (
	INFO     LogLevels = "info"
	WARN     LogLevels = "warn"
	ERROR    LogLevels = "error"
	CRITICAL LogLevels = "critical"
)

type IExecutionContext interface {
	Warn(data any)
	Info(data any)
	Error(err error)
	Close()
}

type ExecutionContext struct {
	id     int64
	origin string
	start  time.Time
	end    time.Time
	logger Logger
}

func init() {
	_logger = func(logLevel LogLevels, id string, data any) {
		log.Printf("[%s] [%s, %s]", logLevel, id, data)
	}
	_middleware = make([]func(id string, origin string, logLevel LogLevels, fields map[string]any), 0)
}

func UseInfluxDb(dsn string, authToken string, org string, bucket string) {
	influxDbLogger := NewInfluxDbLogger(dsn, authToken, org, bucket, nil)
	_middleware = append(_middleware, func(id, origin string, logLevel LogLevels, fields map[string]any) {
		influxDbLogger.Log(id, origin, string(logLevel), fields)
	})
}

func UseInfluxDbWithFailover(dsn string, authToken string, org string, bucket string, logFilePath string) error {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	influxDbLogger := NewInfluxDbLogger(dsn, authToken, org, bucket, file)
	_middleware = append(_middleware, func(id, origin string, logLevel LogLevels, fields map[string]any) {
		influxDbLogger.Log(id, origin, string(logLevel), fields)
	})
	return nil
}

func New(origin string, params ...any) IExecutionContext {
	executionContext := &ExecutionContext{}
	executionContext.origin = origin
	executionContext.logger = _logger
	executionContext.Start(params...)
	return executionContext
}

func NewWithLogger(logger Logger) IExecutionContext {
	executionContext := &ExecutionContext{}
	executionContext.logger = logger
	return executionContext
}

func (e *ExecutionContext) Start(params ...any) {
	e.start = time.Now()
	e.id = time.Now().UnixNano()
	info := make(map[string]any)
	info["status"] = "Started"
	info["params"] = params
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, INFO, info)
	}
}

func (e *ExecutionContext) Error(err error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Errored"
	info["error"] = err.Error()
	e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, ERROR, info)
	}
}

func (e *ExecutionContext) Info(data any) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Executing"
	if data != nil {
		info["data"] = data
	}
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, INFO, info)
	}
}

func (e *ExecutionContext) Warn(data any) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Executing"
	if data != nil {
		info["data"] = data
	}
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, WARN, info)
	}
}

func (e *ExecutionContext) Close() {
	e.end = time.Now()
	if recovered := recover(); recovered != nil {
		info := make(map[string]any)
		info["status"] = "Recovered"
		info["error"] = recovered
		e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
		for _, middleware := range _middleware {
			middleware(fmt.Sprintf("%d", e.id), e.origin, CRITICAL, info)
		}
	}
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Ended"
	info["benchmark"] = e.end.Sub(e.start).Nanoseconds()
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, INFO, info)
	}
}

func RegisterLogger(logger Logger) {
	_logger = logger
}
