package proxy

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	codecs "github.com/vedadiyan/goal/pkg/bus/nats"
	"github.com/vedadiyan/goal/pkg/di"
	"google.golang.org/protobuf/proto"
)

type NATSProxy[TResponse proto.Message] struct {
	conn      *nats.Conn
	codec     codecs.CompressedProtoConn
	namespace string
	new       func() TResponse
}

func (p NATSProxy[TResponse]) Send(request proto.Message) (*TResponse, error) {
	enc, err := p.codec.Encode(p.namespace, request)
	if err != nil {
		return nil, err
	}
	msg, err := p.conn.Request(p.namespace, enc, time.Hour)
	if err != nil {
		return nil, err
	}
	status := msg.Header.Get("status")
	if status != "SUCCESS" {
		return nil, fmt.Errorf(status)
	}
	res := p.new()
	err = p.codec.Decode(p.namespace, msg.Data, res)
	if err != nil {
		return nil, err
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
