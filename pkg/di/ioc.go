package di

import (
	"reflect"
	"sync"
	"time"
)

type Events int

type LifeCycles int

const (
	SINGLETON LifeCycles = iota
	TRANSIENT
	SCOPED
)

const (
	REFRESHING Events = iota
	REFRESHED
)

var (
	_singletonRefresh sync.Map
	_contextTypes     sync.Map
	_context          sync.Map
	_scopedContext    sync.Map
)

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
	return AddSinletonWithName(name, service)
}

func AddSinletonWithName[T any](name string, service func() (instance T, err error)) error {
	singleton := singleton[T]{
		ig: service,
	}
	if _, ok := _context.LoadOrStore(name, &singleton); ok {
		return objectAlreadyExistsError(name)
	}
	_contextTypes.Store(name, SINGLETON)
	return nil
}

func RefreshSinleton[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	return RefreshSinletonWithName(name, service)
}

func RefreshSinletonWithName[T any](name string, service func() (instance T, err error)) error {
	values, ok := _singletonRefresh.Load(name)
	if ok {
		for _, value := range values.([]func(event Events)) {
			value(REFRESHING)
		}
	}
	singleton := singleton[T]{
		ig: service,
	}
	_context.Store(name, &singleton)
	if ok {
		for _, value := range values.([]func(event Events)) {
			value(REFRESHED)
		}
	}
	return nil
}

func OnSingletonRefresh[T any](cb func(event Events)) {
	name := nameOf[T]()
	value, ok := _singletonRefresh.Load(name)
	if !ok {
		value = make([]func(event Events), 0)
	}
	value = append(value.([]func(event Events)), cb)
	_singletonRefresh.Store(name, value)
}

func OnSingletonRefreshWithName(name string, cb func()) {
	value, ok := _singletonRefresh.Load(name)
	if !ok {
		value = make([]func(), 0)
	}
	value = append(value.([]func()), cb)
	_singletonRefresh.Store(name, value)
}

func AddTransient[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	return AddTransientWithName(name, service)
}

func AddTransientWithName[T any](name string, service func() (instance T, err error)) error {
	if _, ok := _context.Load(name); ok {
		return objectAlreadyExistsError(name)
	}
	_context.Store(name, service)
	_contextTypes.Store(name, TRANSIENT)
	return nil
}

func AddScoped[T any](service func() (instance T, err error)) error {
	name := nameOf[T]()
	return AddScopedWithName(name, service)
}

func AddScopedWithName[T any](name string, service func() (instance T, err error)) error {
	if _, ok := _context.LoadOrStore(name, service); ok {
		return objectAlreadyExistsError(name)
	}
	_contextTypes.Store(name, SCOPED)
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
	return ResolveWithName[T](name, options)
}

func ResolveWithName[T any](name string, options *options) (instance *T, err error) {
	lifeCycle, ok := _contextTypes.Load(name)
	if !ok {
		return nil, objectNotFoundError(name)
	}
	object, ok := _context.Load(name)
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
	scopedValue, ok := _scopedContext.Load(options.scopeId)
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
	_scopedContext.Store(options.scopeId, &inst)
	time.AfterFunc(options.ttl, func() {
		_scopedContext.Delete(options.scopeId)
	})
	return &inst, err
}

func nameOf[T any]() string {
	var typeOfT *T
	return reflect.TypeOf(typeOfT).Elem().String()
}

func init() {
	_contextTypes = sync.Map{}
	_context = sync.Map{}
	_scopedContext = sync.Map{}
}
