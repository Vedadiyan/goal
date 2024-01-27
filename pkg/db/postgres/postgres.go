package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vedadiyan/goal/pkg/db/postgres/sanitize"
	"github.com/vedadiyan/goal/pkg/di"
)

type Type int

const (
	COMMAND Type = iota
	QUERY
	COMMAND_TEMPLATE
	QUERY_TEMPLATE
)

type Pool struct {
	pool       *pgxpool.Pool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func getPgxSql(sql string, arguments map[string]any) (string, []any) {
	_sql := sql
	index := 0
	_arguments := make([]any, 0)
	for key, value := range arguments {
		if strings.Contains(_sql, fmt.Sprintf("\"%s\"", key)) {
			_sql = strings.ReplaceAll(_sql, fmt.Sprintf("\"%s\"", key), fmt.Sprintf("$%d", index+1))
			_arguments = append(_arguments, value)
			index++
		}
	}
	return _sql, _arguments
}

func (pool *Pool) Exec(ctx context.Context, sql string, arguments map[string]any) (pgconn.CommandTag, error) {
	str := sql
	if arguments != nil {
		_sql, _arguments := getPgxSql(sql, arguments)
		_str, err := sanitize.SanitizeSQL(_sql, _arguments...)
		if err != nil {
			return pgconn.CommandTag{}, err
		}
		str = _str
	}
	return pool.pool.Exec(ctx, str)
}

func (pool *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return pool.pool.BeginTx(ctx, txOptions)
}

func (pool *Pool) Query(ctx context.Context, sql string, arguments map[string]any) ([]map[string]any, error) {
	str := sql
	if arguments != nil {
		_sql, _arguments := getPgxSql(sql, arguments)
		_str, err := sanitize.SanitizeSQL(_sql, _arguments...)
		if err != nil {
			return nil, err
		}
		str = _str
	}
	res, err := pool.pool.Query(ctx, str)
	if err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0)
	for res.Next() {
		fields := res.FieldDescriptions()
		value, err := res.Values()
		if err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i := 0; i < len(fields); i++ {
			row[fields[i].Name] = value[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
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

func Handle(dsn string, _type Type, sql string, arguments map[string]any) ([]map[string]any, error) {
	pool := di.ResolveWithNameOrPanic[Pool](dsn, nil)
	switch _type {
	case COMMAND:
		{
			res, err := pool.Exec(context.TODO(), sql, arguments)
			if err != nil {
				return nil, err
			}
			return []map[string]any{
				{
					"rows_affected": res.RowsAffected(),
				},
			}, nil
		}
	case QUERY:
		{
			res, err := pool.Query(context.TODO(), sql, arguments)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	case COMMAND_TEMPLATE:
		{
			cmd, err := Build(sql, arguments)
			if err != nil {
				return nil, err
			}
			res, err := pool.Exec(context.TODO(), cmd, nil)
			if err != nil {
				return nil, err
			}
			return []map[string]any{
				{
					"rows_affected": res.RowsAffected(),
				},
			}, nil
		}
	case QUERY_TEMPLATE:
		{
			cmd, err := Build(sql, arguments)
			if err != nil {
				return nil, err
			}
			res, err := pool.Query(context.TODO(), cmd, make(map[string]any))
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	return nil, fmt.Errorf("unsupported operation")
}
