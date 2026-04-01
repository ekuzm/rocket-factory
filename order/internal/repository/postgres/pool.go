package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
)

type Pool struct {
	*pgxpool.Pool
}

type Sqlizer interface {
	ToSql() (string, []any, error)
}

func (p *Pool) Get(ctx context.Context, dst any, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		slog.Error("sql build failed", "operation", "get", "error", err)

		return fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	slog.Debug("sql execute", "operation", "get", "argCount", len(args), "hasTransaction", tx != nil)

	if tx != nil {
		return pgxscan.Get(ctx, tx, dst, query, args...)
	}

	return pgxscan.Get(ctx, p.Pool, dst, query, args...)
}

func (p *Pool) Select(ctx context.Context, dst any, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		slog.Error("sql build failed", "operation", "select", "error", err)

		return fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	slog.Debug("sql execute", "operation", "select", "argCount", len(args), "hasTransaction", tx != nil)

	if tx != nil {
		return pgxscan.Select(ctx, tx, dst, query, args...)
	}

	return pgxscan.Select(ctx, p.Pool, dst, query, args...)
}

func (p *Pool) Exec(ctx context.Context, sqlizer Sqlizer) (pgconn.CommandTag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		slog.Error("sql build failed", "operation", "exec", "error", err)

		return pgconn.CommandTag{}, fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	slog.Debug("sql execute", "operation", "exec", "argCount", len(args), "hasTransaction", tx != nil)

	if tx != nil {
		return tx.Exec(ctx, query, args...)
	}

	return p.Pool.Exec(ctx, query, args...)
}
