package postgres

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	//go:embed stats.table.sql
	_statsTable string
	//go:embed stats.insert.sql
	_statsInsert string
)

type Listener struct {
	conn        *pgx.Conn
	channel     string
	subscribers map[string]func(payload map[string]any)
	mut         sync.Mutex
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

type Msg struct {
	Packet *pgconn.Notification
	Err    error
}

type MsgFrame struct {
	Subject string
	Data    map[string]any
}

func (listener *Listener) next(ctx context.Context) chan *Msg {
	chn := make(chan *Msg, 1)
	packet, err := listener.conn.WaitForNotification(ctx)
	msg := &Msg{
		Packet: packet,
		Err:    err,
	}
	chn <- msg
	return chn
}

func (listener *Listener) tryEnter(ctx context.Context, message string) (bool, error) {
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(message))
	if err != nil {
		return false, err
	}
	bytes := sha256.Sum(nil)
	res, err := listener.conn.Exec(ctx, _statsInsert, hex.EncodeToString(bytes))
	return res.RowsAffected() > 0, err
}

func (listener *Listener) init(ctx context.Context) error {
	_, err := listener.conn.Exec(ctx, _statsTable)
	return err
}

func (listener *Listener) listen(ctx context.Context) error {
	listener.init(ctx)
	cmd := fmt.Sprintf("LISTEN %s;", listener.channel)
	_, err := listener.conn.Exec(ctx, cmd)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case notification := <-listener.next(ctx):
				{
					if notification.Err != nil {
						continue
					}
					check, err := listener.tryEnter(ctx, notification.Packet.Payload)
					if err != nil {
						continue
					}
					if !check {
						continue
					}
					var msgFrame MsgFrame
					err = json.Unmarshal([]byte(notification.Packet.Payload), &msgFrame)
					if err != nil {
						continue
					}
					if handler, ok := listener.subscribers["*"]; ok {
						go handler(map[string]any{"msg": msgFrame})
					}
					if handler, ok := listener.subscribers[msgFrame.Subject]; ok {
						go handler(msgFrame.Data)
					}
				}
			}
		}
	}()
	return nil
}

func (listener *Listener) Subscribe(subject string, handler func(payload map[string]any)) {
	listener.mut.Lock()
	defer listener.mut.Unlock()
	listener.subscribers[subject] = handler
}

func (listener *Listener) Unsubscribe(subject string) {
	listener.mut.Lock()
	defer listener.mut.Unlock()
	delete(listener.subscribers, subject)
}

func (listerner *Listener) Drain() {
	listerner.cancelFunc()
}

func Connect(dsn string, channel string) (*Listener, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cn := &Listener{
		conn:        conn,
		channel:     channel,
		subscribers: make(map[string]func(payload map[string]any)),
		ctx:         ctx,
		cancelFunc:  cancelFunc,
	}
	err = cn.listen(ctx)
	if err != nil {
		return nil, err
	}
	return cn, nil
}
