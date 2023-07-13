package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool       *pgxpool.Pool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func getPgxSql(sql string, arguments map[string]any) (string, []any) {
	_sql := sql
	index := 0
	_arguments := make([]any, len(arguments))
	for key, value := range arguments {
		_sql = strings.ReplaceAll(_sql, key, fmt.Sprintf("$%d", index))
		_arguments[index] = value
		index++
	}
	return _sql, _arguments
}

func (pool *Pool) Exec(ctx context.Context, sql string, arguments map[string]any) (pgconn.CommandTag, error) {
	_sql, _arguments := getPgxSql(sql, arguments)
	return pool.pool.Exec(ctx, _sql, _arguments...)
}

func (pool *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return pool.pool.BeginTx(ctx, txOptions)
}

func (pool *Pool) Query(ctx context.Context, sql string, arguments map[string]any) (pgx.Rows, error) {
	_sql, _arguments := getPgxSql(sql, arguments)
	return pool.pool.Query(ctx, _sql, _arguments...)
}

func (pool *Pool) Close() {
	pool.cancelFunc()
	pool.pool.Close()
}

func New(dsn string, maxConn int, minConn int) (*Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = int32(maxConn)
	conf.MinConns = int32(minConn)
	ctx, cancelFunc := context.WithCancel(context.Background())
	pgpool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	pool := &Pool{
		pool:       pgpool,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	return pool, nil
}
