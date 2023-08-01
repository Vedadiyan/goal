package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"github.com/vedadiyan/goal/pkg/insight"
	"google.golang.org/protobuf/proto"
)

type Option func(*NATSServiceOptions)

type NATSServiceOptions struct {
	isCached  bool
	ttl       time.Duration
	onsuccess []string
	onerror   []string
}

type NATSService[TReq proto.Message, TRes proto.Message, TFuncType ~func(TReq) (TRes, error)] struct {
	conn         *nats.Conn
	codec        *codecs.CompressedProtoConn
	reloadState  chan ReloadStates
	subscription *nats.Subscription
	bucket       *nats.KeyValue
	connName     string
	namespace    string
	queue        string
	handlerFn    TFuncType
	options      NATSServiceOptions
	newReq       func() TReq
	newRes       func() TRes
}

func (t *NATSService[TReq, TRes, TFuncType]) Configure(b bool) {
	if !b {
		di.OnRefreshWithName(t.connName, func(e di.Events) {
			t.reloadState <- RELOADING
			if READY == <-t.reloadState {
				t.reloadState <- RELOADED
				return
			}
		})
	}
	t.conn = *di.ResolveWithNameOrPanic[*nats.Conn](t.connName, nil)
}
func (t *NATSService[TReq, TRes, TFuncType]) configureCache() error {
	js, err := t.conn.JetStream()
	if err != nil {
		return err
	}
	buckets := js.KeyValueStoreNames()
	bucketExists := false
	bucketName := strings.ReplaceAll(t.namespace, ".", "_")
	for bucket := range buckets {
		if bucket == fmt.Sprintf("KV_%s", bucketName) {
			bucketExists = true
			break
		}
	}
	if !bucketExists {
		bucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: bucketName,
			TTL:    t.options.ttl,
		})
		if err != nil {
			return err
		}
		t.bucket = &bucket
		return nil
	}
	bucket, err := js.KeyValue(bucketName)
	if err != nil {
		return err
	}
	t.bucket = &bucket
	return nil
}
func (t *NATSService[TReq, TRes, TFuncType]) Start() error {
	if t.options.isCached {
		err := t.configureCache()
		if err != nil {
			return err
		}
	}
	var subs *nats.Subscription
	var err error
	subs, err = t.conn.QueueSubscribe(t.namespace, t.queue, func(msg *nats.Msg) {
		go t.handler(msg)
	})
	if err != nil {
		return err
	}
	t.subscription = subs
	return nil
}
func (t NATSService[TReq, TRes, TFuncType]) Shutdown() error {
	if t.conn.IsDraining() || t.conn.IsClosed() {
		return nil
	}
	return t.subscription.Unsubscribe()
}
func (t NATSService[TReq, TRes, TFuncType]) Reload() chan ReloadStates {
	return t.reloadState
}
func (t NATSService[TReq, TRes, TFuncType]) error(insight insight.IExecutionContext, originalMsg *nats.Msg, responseMsg *nats.Msg) {
	err := originalMsg.RespondMsg(responseMsg)
	if err != nil {
		insight.Error(err)
	}
	if t.options.onerror == nil {
		return
	}
	onErrorResponse := *responseMsg
	onErrorResponse.Data = originalMsg.Data
	for _, namespace := range t.options.onerror {
		msg := onErrorResponse
		msg.Subject = namespace
		err := t.conn.PublishMsg(&msg)
		if err != nil {
			insight.Error(err)
		}
	}
}
func (t NATSService[TReq, TRes, TFuncType]) success(insight insight.IExecutionContext, originalMsg *nats.Msg, responseMsg *nats.Msg) {
	err := originalMsg.RespondMsg(responseMsg)
	if err != nil {
		insight.Error(err)
	}
	if t.options.onsuccess == nil {
		return
	}
	onSuccessResponse := *responseMsg
	onSuccessResponse.Data = originalMsg.Data
	for _, namespace := range t.options.onsuccess {
		msg := onSuccessResponse
		msg.Subject = namespace
		err := t.conn.PublishMsg(&msg)
		if err != nil {
			insight.Error(err)
		}
	}
}
func (t NATSService[TReq, TRes, TFuncType]) handler(msg *nats.Msg) {
	var requestHash string
	headers := nats.Header{}
	outMsg := &nats.Msg{Subject: msg.Reply, Header: headers}
	insight := insight.New(t.namespace, msg.Reply)
	request := t.newReq()
	insight.OnFailure(func(err error) {
		headers.Add("status", "FAIL:RECOVERED")
		t.error(insight, msg, outMsg)
	})
	insight.Start(request)
	defer insight.Close()
	if t.options.isCached {
		_requestHash, err := GetHash(msg.Data)
		if err != nil {
			insight.Error(err)
			headers.Add("status", "FAIL:REQUEST:HASH")
			t.error(insight, msg, outMsg)
			return
		}
		requestHash = _requestHash
		value, err := (*t.bucket).Get(requestHash)
		if err == nil {
			headers.Add("status", "SUCCESS")
			outMsg.Data = value.Value()
			t.success(insight, msg, outMsg)
			return
		}
	}
	if len(msg.Data) > 0 {
		err := t.codec.Decode(msg.Subject, msg.Data, request)
		if err != nil {
			insight.Error(err)
			headers.Add("status", "FAIL:DECODE")
			t.error(insight, msg, outMsg)
			return
		}
	}
	response, err := t.handlerFn(request)
	if err != nil {
		insight.Error(err)
		headers.Add("status", "FAIL:HANDLE")
		headers.Add("error", err.Error())
		t.error(insight, msg, outMsg)
		return
	}
	bytes, err := t.codec.Encode(msg.Subject, response)
	if err != nil {
		insight.Error(err)
		headers.Add("status", "FAIL:ENCODE")
		t.error(insight, msg, outMsg)
		return
	}
	if t.options.isCached {
		_, err = (*t.bucket).Create(requestHash, bytes)
		if err != nil {
			insight.Warn(err)
		}
	}
	headers.Add("status", "SUCCESS")
	outMsg.Data = bytes
	t.success(insight, msg, outMsg)
}

func GetHash(bytes []byte) (string, error) {
	sha256 := sha256.New()
	_, err := sha256.Write(bytes)
	if err != nil {
		return "", err
	}
	requestHash := sha256.Sum(nil)
	return base64.URLEncoding.EncodeToString(requestHash), nil
}

func New[TReq proto.Message, TRes proto.Message, TFuncType ~func(TReq) (TRes, error)](connName string, namespace string, queue string, handlerFn TFuncType, options ...Option) *NATSService[TReq, TRes, TFuncType] {
	tReq := reflect.TypeOf(*new(TReq)).Elem()
	tRes := reflect.TypeOf(*new(TRes)).Elem()
	service := NATSService[TReq, TRes, TFuncType]{
		namespace:   namespace,
		queue:       queue,
		handlerFn:   handlerFn,
		connName:    connName,
		reloadState: make(chan ReloadStates),
		newReq: func() TReq {
			return reflect.New(tReq).Interface().(TReq)
		},
		newRes: func() TRes {
			return reflect.New(tRes).Interface().(TRes)
		},
		codec: &codecs.CompressedProtoConn{},
	}
	for _, option := range options {
		option(&service.options)
	}
	return &service
}

func WithCache(ttl time.Duration) Option {
	return func(n *NATSServiceOptions) {
		n.isCached = true
		n.ttl = ttl
	}
}

func WithOnSuccessCallBacks(namespaces ...string) Option {
	return func(no *NATSServiceOptions) {
		no.onsuccess = make([]string, 0)
		no.onsuccess = append(no.onsuccess, namespaces...)
	}
}

func WithOnFailureCallBacks(namespaces ...string) Option {
	return func(no *NATSServiceOptions) {
		no.onerror = make([]string, 0)
		no.onerror = append(no.onerror, namespaces...)
	}
}
