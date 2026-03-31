package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
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
		logger.WithFields(logrus.Fields{
			"Operation": "Get",
			"error":     err,
		}).Error("Failed to build SQL query")

		return fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	logger.WithFields(logrus.Fields{
		"Operation":       "Get",
		"Query":           query,
		"Args":            args,
		"Has Transaction": tx != nil,
	}).Debug("Execute SQL query")

	if tx != nil {
		return pgxscan.Get(ctx, tx, dst, query, args...)
	}

	return pgxscan.Get(ctx, p.Pool, dst, query, args...)
}

func (p *Pool) Select(ctx context.Context, dst any, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Operation": "Select",
			"error":     err,
		}).Error("Failed to build SQL query")

		return fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	logger.WithFields(logrus.Fields{
		"Operation":       "Select",
		"Query":           query,
		"Args":            args,
		"Has Transaction": tx != nil,
	}).Debug("Execute SQL query")

	if tx != nil {
		return pgxscan.Select(ctx, tx, dst, query, args...)
	}

	return pgxscan.Select(ctx, p.Pool, dst, query, args...)
}

func (p *Pool) Exec(ctx context.Context, sqlizer Sqlizer) (pgconn.CommandTag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Operation": "Exec",
			"error":     err,
		}).Error("Failed to build SQL query")

		return pgconn.CommandTag{}, fmt.Errorf("sqlizer.ToSql: %w", err)
	}

	tx := transaction.Extract(ctx)
	logger.WithFields(logrus.Fields{
		"Operation":       "Exec",
		"Query":           query,
		"Args":            args,
		"Has Transaction": tx != nil,
	}).Debug("Execute SQL query")

	if tx != nil {
		return tx.Exec(ctx, query, args...)
	}

	return p.Pool.Exec(ctx, query, args...)
}
