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

type IExecutionContext[T any] interface {
	Start(params ...any)
	Interrupt(err *T) (*T, error)
	Error(err error) (*T, error)
	Complete(value *T) (*T, error)
	Recover()
}

type ExecutionContext[T any] struct {
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

func New[T any](origin string) IExecutionContext[T] {
	executionContext := &ExecutionContext[T]{}
	executionContext.origin = origin
	executionContext.logger = _logger
	return executionContext
}

func NewWithLogger[T any](logger Logger) IExecutionContext[T] {
	executionContext := &ExecutionContext[T]{}
	executionContext.logger = logger
	return executionContext
}

func (e *ExecutionContext[T]) Start(params ...any) {
	e.start = time.Now()
	e.id = time.Now().UnixNano()
	info := make(map[string]any)
	info["status"] = "Executing"
	info["params"] = params
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, INFO, info)
	}
}

func (e *ExecutionContext[T]) Interrupt(err *T) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Interuppted"
	info["reason"] = err
	e.logger(WARN, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, WARN, info)
	}
	return err, nil
}

func (e *ExecutionContext[T]) Error(err error) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Errored"
	info["error"] = err.Error()
	e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, ERROR, info)
	}
	return nil, err
}

func (e *ExecutionContext[T]) Complete(value *T) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Executed"
	info["benchmark"] = e.end.Sub(e.start).Nanoseconds()
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	for _, middleware := range _middleware {
		middleware(fmt.Sprintf("%d", e.id), e.origin, INFO, info)
	}
	return value, nil
}

func (e *ExecutionContext[T]) Recover() {
	if recovered := recover(); recovered != nil {
		e.end = time.Now()
		info := make(map[string]any)
		info["status"] = "Recovered"
		info["error"] = recovered
		e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
		for _, middleware := range _middleware {
			middleware(fmt.Sprintf("%d", e.id), e.origin, CRITICAL, info)
		}
	}
}

func RegisterLogger(logger Logger) {
	_logger = logger
}
