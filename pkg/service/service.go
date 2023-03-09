package service

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"github.com/vedadiyan/goal/pkg/insight"
	"google.golang.org/protobuf/proto"
)

type Option func(*NATSService)

type Handler func(proto.Message) (proto.Message, error)

type NATSService struct {
	conn         *nats.Conn
	codec        *codecs.CompressedProtoConn
	reloadState  chan ReloadStates
	subscription *nats.Subscription
	bucket       *nats.KeyValue
	isCached     bool
	ttl          time.Duration
	connName     string
	namespace    string
	queue        string
	handlerFn    Handler
}

func (t *NATSService) Configure(b bool) {
	if !b {
		di.OnSingletonRefreshWithName(t.connName, func(e di.Events) {
			if e == di.REFRESHED {
				t.conn = *di.ResolveWithNameOrPanic[*nats.Conn](t.connName, nil)
				t.reloadState <- RELOADED
				return
			}
			t.reloadState <- RELOADING
		})
		return
	}
	t.conn = *di.ResolveWithNameOrPanic[*nats.Conn](t.connName, nil)
}
func (t *NATSService) configureCache() error {
	js, err := t.conn.JetStream()
	if err != nil {
		return err
	}
	buckets := js.KeyValueStoreNames()
	bucketExists := false
	for bucket := range buckets {
		if bucket == t.namespace {
			bucketExists = true
			break
		}
	}
	if !bucketExists {
		bucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: t.namespace,
			TTL:    t.ttl,
		})
		if err != nil {
			return err
		}
		t.bucket = &bucket
		return nil
	}
	bucket, err := js.KeyValue(t.namespace)
	if err != nil {
		return err
	}
	t.bucket = &bucket
	return nil
}
func (t *NATSService) Start() error {
	if t.isCached {
		err := t.configureCache()
		if err != nil {
			return err
		}
	}
	subs, err := t.conn.QueueSubscribe(t.namespace, t.queue, func(msg *nats.Msg) {
		go t.handler(msg)
	})
	if err != nil {
		return err
	}
	t.subscription = subs
	return nil
}
func (t NATSService) Shutdown() error {
	return t.subscription.Unsubscribe()
}
func (t NATSService) Reload() <-chan ReloadStates {
	return t.reloadState
}
func (t NATSService) handler(msg *nats.Msg) {
	insight := insight.New(t.namespace, msg.Reply)
	defer insight.Close()
	var request proto.Message
	headers := nats.Header{}
	outMsg := &nats.Msg{Subject: msg.Reply, Header: headers}
	if t.isCached {
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
	if t.isCached {
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

func New(connName string, namespace string, queue string, handlerFn Handler, options ...Option) *NATSService {
	service := NATSService{
		namespace: namespace,
		queue:     queue,
		handlerFn: handlerFn,
		connName:  connName,
	}
	for _, option := range options {
		option(&service)
	}
	return &service
}

func WithCache(ttl time.Duration) Option {
	return func(n *NATSService) {
		n.isCached = true
		n.ttl = ttl
	}
}
