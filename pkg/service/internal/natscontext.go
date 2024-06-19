package natsCtx

import (
	"github.com/nats-io/nats.go"
	"github.com/vedadiyan/goal/pkg/insight"
)

type Header map[string]string

type NatsCtx struct {
	conn        *nats.Conn
	insight     insight.IExecutionContext
	requestMsg  *nats.Msg
	responseMsg *nats.Msg
}

func NewNatsCtx(conn *nats.Conn, insight insight.IExecutionContext, msg *nats.Msg) *NatsCtx {
	return &NatsCtx{
		conn:        conn,
		insight:     insight,
		requestMsg:  msg,
		responseMsg: &nats.Msg{Subject: msg.Reply, Header: nats.Header{}},
	}
}

func (nc *NatsCtx) Error(headers Header) {
	msg := &nats.Msg{}
	msg.Header = nats.Header{}
	for key, value := range headers {
		msg.Header.Add(key, value)
	}
	err := nc.requestMsg.RespondMsg(msg)
	if err != nil {
		nc.insight.Error(err)
	}

}
func (nc *NatsCtx) Success(data []byte, headers Header) {
	msg := &nats.Msg{}
	msg.Header = nats.Header{}
	for key, value := range headers {
		msg.Header.Add(key, value)
	}
	msg.Data = data
	err := nc.requestMsg.RespondMsg(msg)
	if err != nil {
		nc.insight.Error(err)
	}
}
