package proxy

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"google.golang.org/protobuf/proto"
)

type ProxyError string

const (
	_ENCODE_ERROR  ProxyError = ProxyError("encode error")
	_DECODE_ERROR  ProxyError = ProxyError("decode error")
	_GATEWAY_ERROR ProxyError = ProxyError("gateway error")
)

func (p ProxyError) Error() string {
	return string(p)
}

type NATSProxy[TResponse proto.Message] struct {
	conn      *nats.Conn
	codec     codecs.CompressedProtoConn
	namespace string
	new       func() TResponse
}

func (p NATSProxy[TResponse]) Send(request proto.Message) (*TResponse, error) {
	enc, err := p.codec.Encode(p.namespace, request)
	if err != nil {
		return nil, _ENCODE_ERROR
	}
	msg, err := p.conn.Request(p.namespace, enc, time.Hour)
	if err != nil {
		return nil, _GATEWAY_ERROR

	}
	status := msg.Header.Get("status")
	if status != "SUCCESS" {
		return nil, fmt.Errorf(`{"status": "%s", "message": "%s"}`, status, string(msg.Data))
	}
	res := p.new()
	err = p.codec.Decode(p.namespace, msg.Data, res)
	if err != nil {
		return nil, _DECODE_ERROR
	}
	return &res, nil
}
func New[TResponse proto.Message](connName string, namespace string, newRes func() TResponse) *NATSProxy[TResponse] {
	conn := *di.ResolveWithNameOrPanic[*nats.Conn](connName, nil)
	natsProxy := NATSProxy[TResponse]{
		namespace: namespace,
		conn:      conn,
		codec:     codecs.CompressedProtoConn{},
		new:       newRes,
	}
	di.OnRefreshWithName(connName, func(e di.Events) {
		natsProxy.conn = *di.ResolveWithNameOrPanic[*nats.Conn](connName, nil)
	})
	return &natsProxy
}

func Create[TResponse any](connName string, namespace string) *NATSProxy[proto.Message] {
	conn := *di.ResolveWithNameOrPanic[*nats.Conn](connName, nil)
	natsProxy := NATSProxy[proto.Message]{
		namespace: namespace,
		conn:      conn,
		codec:     codecs.CompressedProtoConn{},
		new: func() proto.Message {
			var res TResponse
			return any(&res).(proto.Message)
		},
	}
	di.OnRefreshWithName(connName, func(e di.Events) {
		natsProxy.conn = *di.ResolveWithNameOrPanic[*nats.Conn](connName, nil)
	})
	return &natsProxy
}
