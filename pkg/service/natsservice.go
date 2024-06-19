package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"github.com/vedadiyan/goal/pkg/insight"
	internal "github.com/vedadiyan/goal/pkg/service/internal"
	"google.golang.org/protobuf/proto"
)

type NATSService[TReq proto.Message, TRes proto.Message, TFuncType ~func(TReq) (int, TRes, error)] struct {
	conn         *nats.Conn
	codec        *codecs.CompressedProtoConn
	reloadState  chan ReloadStates
	subscription *nats.Subscription
	connName     string
	namespace    string
	queue        string
	handlerFn    TFuncType
	newReq       func() TReq
	newRes       func() TRes
}

const (
	HEADER_STATUS = "X-Status"
	HEADER_ERROR  = "X-Error"

	STATUS_SERVER_FAILURE  = "500"
	STATUS_SERVICE_FAILURE = "502"
)

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
	t.conn = di.ResolveWithNameOrPanic[nats.Conn](t.connName, nil)
}
func (t *NATSService[TReq, TRes, TFuncType]) Start() error {
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

func (t NATSService[TReq, TRes, TFuncType]) handler(msg *nats.Msg) {
	insight := insight.New(t.namespace, msg.Reply)
	defer insight.Close()
	ctx := internal.NewNatsCtx(t.conn, insight, msg)
	request := t.newReq()
	insight.OnFailure(func(err error) {
		ctx.Error(internal.Header{HEADER_STATUS: STATUS_SERVER_FAILURE, HEADER_ERROR: err.Error()})
	})
	if len(msg.Data) > 0 {
		err := t.codec.Decode(msg.Subject, msg.Data, request)
		if err != nil {
			insight.Error(err)
			ctx.Error(internal.Header{HEADER_STATUS: STATUS_SERVER_FAILURE, HEADER_ERROR: err.Error()})
			return
		}
	}
	insight.Start(request)
	status, response, err := t.handlerFn(request)
	if err != nil {
		insight.Error(err)
		if status == 0 {
			ctx.Error(internal.Header{HEADER_STATUS: STATUS_SERVICE_FAILURE, HEADER_ERROR: strings.ReplaceAll(err.Error(), "\"", "\\\"")})
			return
		}
		ctx.Error(internal.Header{HEADER_STATUS: fmt.Sprintf("%v", status), HEADER_ERROR: strings.ReplaceAll(err.Error(), "\"", "\\\"")})
		return
	}
	bytes, err := t.codec.Encode(msg.Subject, response)
	if err != nil {
		insight.Error(err)
		ctx.Error(internal.Header{HEADER_STATUS: STATUS_SERVER_FAILURE, HEADER_ERROR: err.Error()})
		return
	}
	ctx.Success(bytes, internal.Header{HEADER_STATUS: fmt.Sprintf("%v", status)})
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

func New[TReq proto.Message, TRes proto.Message, TFuncType ~func(TReq) (int, TRes, error)](connName string, namespace string, queue string, handlerFn TFuncType) *NATSService[TReq, TRes, TFuncType] {
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
	return &service
}
