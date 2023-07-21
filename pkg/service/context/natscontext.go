package context

import (
	"github.com/nats-io/nats.go"
	"github.com/vedadiyan/goal/pkg/insight"
)

type NATSContext struct {
	conn        *nats.Conn
	insight     insight.IExecutionContext
	requestMsg  *nats.Msg
	responseMsg *nats.Msg
	onsuccess   []string
	onerror     []string
}

func NewNatsContext(conn *nats.Conn, insight insight.IExecutionContext, msg *nats.Msg, onerror []string, onsuccess []string) *NATSContext {
	return &NATSContext{
		conn:        conn,
		insight:     insight,
		requestMsg:  msg,
		responseMsg: &nats.Msg{Subject: msg.Reply, Header: nats.Header{}},
		onsuccess:   onsuccess,
		onerror:     onerror,
	}
}

func (nc *NATSContext) Error(headers map[string]string) {
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
func (nc *NATSContext) Success(data []byte, headers map[string]string) {
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
