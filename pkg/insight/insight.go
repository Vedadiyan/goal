package insight

import (
	"fmt"
	"log"
	"time"
)

type Logger func(logLevel LogLevels, id string, data any)

type LogLevels string

var (
	_logger Logger
)

const (
	INFO  LogLevels = "info"
	WARN  LogLevels = "warn"
	ERROR LogLevels = "error"
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
	start  time.Time
	end    time.Time
	logger Logger
}

func init() {
	_logger = func(logLevel LogLevels, id string, data any) {
		log.Printf("[%s] [%s, %s]", logLevel, id, data)
	}
}

func New[T any]() IExecutionContext[T] {
	executionContext := &ExecutionContext[T]{}
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
}

func (e *ExecutionContext[T]) Interrupt(err *T) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Interuppted"
	info["time"] = e.end.Sub(e.start).Nanoseconds()
	info["reason"] = err
	e.logger(WARN, fmt.Sprintf("%d", e.id), info)
	return err, nil
}

func (e *ExecutionContext[T]) Error(err error) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Errored"
	info["time"] = e.end.Sub(e.start).Nanoseconds()
	info["error"] = err.Error()
	e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
	return nil, err
}

func (e *ExecutionContext[T]) Complete(value *T) (*T, error) {
	e.end = time.Now()
	info := make(map[string]any)
	info["status"] = "Executed"
	info["time"] = e.end.Sub(e.start).Nanoseconds()
	e.logger(INFO, fmt.Sprintf("%d", e.id), info)
	return value, nil
}

func (e *ExecutionContext[T]) Recover() {
	if recovered := recover(); recovered != nil {
		e.end = time.Now()
		info := make(map[string]any)
		info["status"] = "Recovered"
		info["time"] = e.end.Sub(e.start).Nanoseconds()
		info["error"] = recovered
		e.logger(ERROR, fmt.Sprintf("%d", e.id), info)
	}
}

func RegisterLogger(logger Logger) {
	_logger = logger
}
