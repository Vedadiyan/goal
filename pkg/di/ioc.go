package di

import (
	"reflect"
	"sync"
	"time"
)

type LifeCycles int

const (
	SINGLETON LifeCycles = iota
	TRANSIENT
	SCOPED
)

var contextTypes sync.Map
var context sync.Map
var scopedContext sync.Map

type options struct {
	scopeId uint64
	ttl     time.Duration
}

type singleton[T any] struct {
	ig       func() (instance T, err error)
	created  bool
	instance *T
	err      error
	once     sync.Once
}

func (s *singleton[T]) getInstance() (instance T, err error) {
	s.once.Do(func() {
		value, err := s.ig()
		s.instance = &value
		s.err = err
		s.created = true
	})
	return *s.instance, s.err
}

func NewOptions(scopeId uint64, ttl time.Duration) *options {
	return &options{scopeId, ttl}
}

func AddSinleton[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	singleton := singleton[T]{
		ig: service,
	}
	if _, ok := context.LoadOrStore(name, &singleton); ok {
		return objectAlreadyExistsError(name)
	}
	contextTypes.Store(name, SINGLETON)
	return nil
}

func AddTransient[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	if _, ok := context.Load(name); ok {
		return objectAlreadyExistsError(name)
	}
	context.Store(name, service)
	contextTypes.Store(name, TRANSIENT)
	return nil
}

func AddScoped[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	if _, ok := context.LoadOrStore(name, service); ok {
		return objectAlreadyExistsError(name)
	}
	contextTypes.Store(name, SCOPED)
	return nil
}

func ResolveOrPanic[T any](options *options) *T {
	value, err := Resolve[T](options)
	if err != nil {
		panic(err)
	}
	return value
}

func ResolveOrNil[T any](options *options) *T {
	value, _ := Resolve[T](options)
	return value
}

func Resolve[T any](options *options) (instance *T, err error) {
	name := nameOf[T]()
	lifeCycle, ok := contextTypes.Load(name)
	if !ok {
		return nil, objectNotFoundError(name)
	}
	object, ok := context.Load(name)
	if !ok {
		return nil, objectNotFoundError(name)
	}
	switch lifeCycle {
	case SINGLETON:
		return resolveSingleton[T](object, name)
	case TRANSIENT:
		return resolveTransient[T](object, name)
	case SCOPED:
		return resolveScoped[T](options, object, name)
	default:
		return nil, nil
	}
}

func resolveSingleton[T any](object any, name string) (instance *T, err error) {
	value, ok := object.(*singleton[T])
	if !ok {
		return nil, invalidCastError(name)
	}
	inst, err := value.getInstance()
	return &inst, err
}

func resolveTransient[T any](object any, name string) (instance *T, err error) {
	value, ok := object.(func() (instance T, err error))
	if !ok {
		return nil, invalidCastError(name)
	}
	inst, err := value()
	return &inst, err
}

func resolveScoped[T any](options *options, object any, name string) (instance *T, err error) {
	if options == nil {
		return nil, missingRequiredParameter("Options")
	}
	scopedValue, ok := scopedContext.Load(options.scopeId)
	if ok {
		if value, ok := scopedValue.(*T); ok {
			return value, nil
		}
		return nil, invalidCastError(name)
	}
	value, ok := object.(func() (instance T, err error))
	if !ok {
		return nil, invalidCastError(name)
	}
	inst, err := value()
	scopedContext.Store(options.scopeId, &inst)
	time.AfterFunc(options.ttl, func() {
		scopedContext.Delete(options.scopeId)
	})
	return &inst, err
}

func nameOf[T any]() string {
	var typeOfT *T
	return reflect.TypeOf(typeOfT).Elem().String()
}

func init() {
	contextTypes = sync.Map{}
	context = sync.Map{}
	scopedContext = sync.Map{}
}
