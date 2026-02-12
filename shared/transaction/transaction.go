package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Init(p *pgxpool.Pool) {
	pool = p
}

type ctxKey struct{}

type Executor interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func Extract(ctx context.Context) Executor {
	tx, ok := ctx.Value(ctxKey{}).(pgx.Tx)
	if !ok {
		return pool
	}

	return tx
}
