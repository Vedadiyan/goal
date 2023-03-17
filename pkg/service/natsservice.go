package service

import (
	"crypto/sha256"
	"encoding/base64"
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
	isCached bool
	ttl      time.Duration
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
		if bucket == bucketName {
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
func (t NATSService[TReq, TRes, TFuncType]) handler(msg *nats.Msg) {
	insight := insight.New(t.namespace, msg.Reply)
	defer insight.Close()
	var request TReq
	headers := nats.Header{}
	outMsg := &nats.Msg{Subject: msg.Reply, Header: headers}
	if t.options.isCached {
		requestHash, err := GetHash(msg.Data)
		if err != nil {
			headers.Add("status", "FAIL:REQUEST:HASH")
			msg.RespondMsg(outMsg)
			insight.Error(err)
			return
		}
		value, err := (*t.bucket).Get(requestHash)
		if err == nil {
			headers.Add("status", "SUCCESS")
			outMsg.Data = value.Value()
			msg.RespondMsg(outMsg)
			if err != nil {
				insight.Error(err)
				return
			}
			return
		}
	}
	err := t.codec.Decode(msg.Subject, msg.Data, &request)
	if err != nil {
		headers.Add("status", "FAIL:DECODE")
		msg.RespondMsg(outMsg)
		insight.Error(err)
		return
	}
	insight.Start(request)
	response, err := t.handlerFn(request)
	if err != nil {
		headers.Add("status", "FAIL:HANDLE")
		msg.RespondMsg(outMsg)
		insight.Error(err)
		return
	}
	bytes, err := t.codec.Encode(msg.Subject, response)
	if err != nil {
		headers.Add("status", "FAIL:ENCODE")
		msg.RespondMsg(outMsg)
		insight.Error(err)
		return
	}
	if t.options.isCached {
		requestHash, err := GetHash(bytes)
		if err != nil {
			headers.Add("status", "FAIL:REQUEST:HASH")
			msg.RespondMsg(outMsg)
			insight.Error(err)
			return
		}
		(*t.bucket).Create(requestHash, bytes)
	}
	headers.Add("status", "SUCCESS")
	outMsg.Data = bytes
	msg.RespondMsg(outMsg)
	if err != nil {
		insight.Error(err)
		return
	}
}

func GetHash(bytes []byte) (string, error) {
	sha256 := sha256.New()
	_, err := sha256.Write(bytes)
	if err != nil {
		return "", err
	}
	requestHash := sha256.Sum(nil)
	return base64.StdEncoding.EncodeToString(requestHash), nil
}

func New[TReq proto.Message, TRes proto.Message, TFuncType ~func(TReq) (TRes, error)](connName string, namespace string, queue string, handlerFn TFuncType, options ...Option) *NATSService[TReq, TRes, TFuncType] {
	service := NATSService[TReq, TRes, TFuncType]{
		namespace:   namespace,
		queue:       queue,
		handlerFn:   handlerFn,
		connName:    connName,
		reloadState: make(chan ReloadStates),
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
