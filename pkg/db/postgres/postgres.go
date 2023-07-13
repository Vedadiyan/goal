package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool       *pgxpool.Pool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New(dsn string) (*Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
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
