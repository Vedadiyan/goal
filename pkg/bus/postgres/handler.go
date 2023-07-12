package postgres

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
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

type Connection struct {
	conn        *pgx.Conn
	subscribers map[string]func(payload string)
	mut         sync.Mutex
}

type Msg struct {
	Packet *pgconn.Notification
	Err    error
}

func (conn *Connection) next(ctx context.Context) chan *Msg {
	chn := make(chan *Msg)
	packet, err := conn.conn.WaitForNotification(ctx)
	msg := &Msg{
		Packet: packet,
		Err:    err,
	}
	chn <- msg
	return chn
}

func (conn *Connection) check(ctx context.Context, message string) (bool, error) {
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(message))
	if err != nil {
		return false, err
	}
	bytes := sha256.Sum(nil)
	res, err := conn.conn.Exec(ctx, _statsInsert, hex.EncodeToString(bytes))
	return res.RowsAffected() > 0, err
}

func (conn *Connection) init(ctx context.Context) error {
	_, err := conn.conn.Exec(ctx, _statsTable)
	return err
}

func (conn *Connection) Listen(ctx context.Context) {
	conn.init(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case notification := <-conn.next(ctx):
				{
					if notification.Err != nil {
						return
					}
					check, err := conn.check(ctx, notification.Packet.Payload)
					if err != nil {
						return
					}
					if !check {
						return
					}
					if handler, ok := conn.subscribers[notification.Packet.Channel]; ok {
						handler(notification.Packet.Payload)
					}
				}
			}
		}
	}()
}

func (conn *Connection) Subscribe(subject string, handler func(payload string)) {
	conn.mut.Lock()
	defer conn.mut.Unlock()
	conn.subscribers[subject] = handler
}

func (conn *Connection) Unsubscribe(subject string) {
	conn.mut.Lock()
	defer conn.mut.Unlock()
	delete(conn.subscribers, subject)
}
