package transaction

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type manager struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *manager {
	return &manager{
		pool: pool,
	}
}

type Executor interface {
	QueryRow(context.Context, string, ...any) sql.Row
	Query(context.Context, string, ...any) (sql.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type ctxKey struct{}

func inject(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func Extract(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(ctxKey{}).(pgx.Tx)
	if !ok {
		return nil
	}

	return tx
}
