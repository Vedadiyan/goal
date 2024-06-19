package proxy

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"google.golang.org/protobuf/proto"
)

type ProxyError string

const (
	HEADER_STATUS = "X-Status"
	HEADER_ERROR  = "X-Error"

	_ERR_ENCODE    ProxyError = ProxyError("encode error")
	_ERR_DECODE    ProxyError = ProxyError("decode error")
	_ERROR_GATEWAY ProxyError = ProxyError("gateway error")
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
		return nil, _ERR_ENCODE
	}
	msg, err := p.conn.Request(p.namespace, enc, time.Hour)
	if err != nil {
		return nil, _ERROR_GATEWAY

	}
	status := msg.Header.Get(HEADER_STATUS)
	if strings.HasPrefix(status, "5") || strings.HasPrefix(status, "4") {
		return nil, fmt.Errorf(`{"status": "%s", "message": "%s"}`, status, msg.Header.Get(HEADER_ERROR))
	}
	res := p.new()
	err = p.codec.Decode(p.namespace, msg.Data, res)
	if err != nil {
		return nil, _ERR_DECODE
	}
	return &res, nil
}
func New[TResponse proto.Message](connName string, namespace string, newRes func() TResponse) *NATSProxy[TResponse] {
	conn := di.ResolveWithNameOrPanic[nats.Conn](connName, nil)
	natsProxy := NATSProxy[TResponse]{
		namespace: namespace,
		conn:      conn,
		codec:     codecs.CompressedProtoConn{},
		new:       newRes,
	}
	di.OnRefreshWithName(connName, func(e di.Events) {
		natsProxy.conn = di.ResolveWithNameOrPanic[nats.Conn](connName, nil)
	})
	return &natsProxy
}

func Create[TResponse any](connName string, namespace string) *NATSProxy[proto.Message] {
	conn := di.ResolveWithNameOrPanic[nats.Conn](connName, nil)
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
		natsProxy.conn = di.ResolveWithNameOrPanic[nats.Conn](connName, nil)
	})
	return &natsProxy
}
