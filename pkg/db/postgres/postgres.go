package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func getPgxSql(sql string, arguments map[string]any) (string, error) {
	_sql := sql
	for key, value := range arguments {
		if strings.Contains(_sql, fmt.Sprintf("\"$%s\"", key)) {
			v, err := sanitize.SanitizeSQL("$1", standardize(value))
			if err != nil {
				return "", err
			}
			_sql = strings.ReplaceAll(_sql, fmt.Sprintf("\"$%s\"", key), v)
		}
	}
	return _sql, nil
}

func (pool *Pool) Exec(ctx context.Context, sql string, arguments map[string]any) (pgconn.CommandTag, error) {
	str := sql
	if arguments != nil {
		_sql, err := getPgxSql(sql, arguments)
		if err != nil {
			return pgconn.CommandTag{}, err
		}
		str = _sql
	}
	return pool.pool.Exec(ctx, str)
}

func (pool *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return pool.pool.BeginTx(ctx, txOptions)
}

func (pool *Pool) Query(ctx context.Context, sql string, arguments map[string]any) ([]map[string]any, error) {
	str := sql
	if arguments != nil {
		_sql, err := getPgxSql(sql, arguments)
		if err != nil {
			return nil, err
		}
		str = _sql
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
			val := value[i]
			switch t := val.(type) {
			case pgtype.Range[any]:
				{
					out := []any{
						fmt.Sprintf("%v", t.Lower),
						fmt.Sprintf("%v", t.Upper),
					}
					row[fields[i].Name] = out
				}
			default:
				{
					oid := fields[i].DataTypeOID
					if oid == 790 {
						str := fmt.Sprintf("%v", val)
						str = strings.Replace(str, "$", "", 1)
						n, err := strconv.ParseFloat(str, 64)
						if err != nil {
							return nil, err
						}
						row[fields[i].Name] = n
						continue
					}
					if oid == 114 {
						err := Fix(value[i])
						if err != nil {
							return nil, err
						}
					}
					row[fields[i].Name] = value[i]
				}
			}

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
			res, err := pool.Exec(context.TODO(), cmd, arguments)
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
			res, err := pool.Query(context.TODO(), cmd, arguments)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	return nil, fmt.Errorf("unsupported operation")
}

func Unmarshall(str string) (any, error) {
	data := make(map[string]any)
	err := json.Unmarshal([]byte(fmt.Sprintf(`{"root": %s}`, str)), &data)
	if err != nil {
		return nil, err
	}
	err = Fix(data)
	if err != nil {
		return nil, err
	}
	return data["root"], nil
}

// This is a temporary fix for PGX's inability to parse data as expected
// TO DO: Refactor and fix issues
func Fix(data any) error {
	switch t := data.(type) {
	case map[string]any:
		{
			data := t
			for key, value := range data {
				switch t := value.(type) {
				case string:
					{
						if strings.HasPrefix(t, "$") {
							str := strings.TrimPrefix(t, "$")
							n, err := strconv.ParseFloat(str, 64)
							if err != nil {
								return err
							}
							data[key] = n
							continue
						}
						if (strings.HasPrefix(t, "(") || strings.HasPrefix(t, "[")) && (strings.HasSuffix(t, ")") || strings.HasSuffix(t, "]")) {
							segments := strings.Split(t, ",")
							if len(segments) != 2 {
								continue
							}
							left := strings.TrimPrefix(segments[0], "(")
							left = strings.TrimPrefix(left, "[")
							right := strings.TrimRight(segments[1], ")")
							right = strings.TrimRight(right, "]")
							data[key] = []any{left, right}
						}
					}
				case []map[string]any:
					{
						for _, item := range t {
							err := Fix(item)
							if err != nil {
								return err
							}
						}
					}
				case map[string]any:
					{
						err := Fix(t)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	case []any:
		{
			for _, item := range t {
				err := Fix(item)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
