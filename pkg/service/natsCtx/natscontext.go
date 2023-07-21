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
	onsuccess   []string
	onerror     []string
}

func New(conn *nats.Conn, insight insight.IExecutionContext, msg *nats.Msg, onerror []string, onsuccess []string) *NatsCtx {
	return &NatsCtx{
		conn:        conn,
		insight:     insight,
		requestMsg:  msg,
		responseMsg: &nats.Msg{Subject: msg.Reply, Header: nats.Header{}},
		onsuccess:   onsuccess,
		onerror:     onerror,
	}
}

func (nc *NatsCtx) Error(headers Header) {
	for key, value := range headers {
		nc.responseMsg.Header.Add(key, value)
	}
	err := nc.requestMsg.RespondMsg(nc.responseMsg)
	if err != nil {
		nc.insight.Error(err)
	}
	if nc.onerror == nil {
		return
	}
	onErrorResponse := *nc.responseMsg
	onErrorResponse.Reply = ""
	onErrorResponse.Data = nc.requestMsg.Data
	for _, namespace := range nc.onerror {
		msg := onErrorResponse
		msg.Subject = namespace
		err := nc.conn.PublishMsg(&msg)
		if err != nil {
			nc.insight.Error(err)
		}
	}

}
func (nc *NatsCtx) Success(data []byte, headers Header) {
	for key, value := range headers {
		nc.responseMsg.Header.Add(key, value)
	}
	err := nc.requestMsg.RespondMsg(nc.responseMsg)
	if err != nil {
		nc.insight.Error(err)
	}
	if nc.onsuccess == nil {
		return
	}
	onSuccessResponse := *nc.responseMsg
	onSuccessResponse.Reply = ""
	onSuccessResponse.Data = data
	for _, namespace := range nc.onsuccess {
		msg := onSuccessResponse
		msg.Subject = namespace
		err := nc.conn.PublishMsg(&msg)
		if err != nil {
			nc.insight.Error(err)
		}
	}
}
