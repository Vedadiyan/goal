package service

import (
	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"github.com/vedadiyan/goal/pkg/insight"
	"google.golang.org/protobuf/proto"
)

type Handler func(proto.Message) (proto.Message, error)

type NATSService struct {
	conn         *nats.Conn
	codec        *codecs.CompressedProtoConn
	reloadState  chan ReloadStates
	subscription *nats.Subscription

	connName  string
	namespace string
	queue     string
	handlerFn Handler
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
func (t *NATSService) Start() error {
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
	id := msg.Reply
	insight := insight.New(t.namespace, id)
	defer insight.Close()
	var request proto.Message
	headers := nats.Header{}
	outMsg := &nats.Msg{Subject: msg.Reply, Header: headers}
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
	headers.Add("status", "SUCCESS")
	outMsg.Data = bytes
	msg.RespondMsg(outMsg)
	if err != nil {
		insight.Error(err)
		return
	}
}
func New(connName string, namespace string, queue string, handlerFn Handler) *NATSService {
	service := NATSService{
		namespace: namespace,
		queue:     queue,
		handlerFn: handlerFn,
		connName:  connName,
	}
	return &service
}
