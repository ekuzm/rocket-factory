package transaction

import (
	"context"
	"fmt"
	"log/slog"
)

func (m *manager) Wrap(ctx context.Context, callback func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		slog.Error("tx begin failed", "error", err)

		return fmt.Errorf("create transaction: %w", err)
	}

	inject(ctx, tx)
	if err = callback(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			slog.Error("tx rollback failed", "callbackError", err, "error", rollbackErr)

			return fmt.Errorf("rollback transaction: %w", rollbackErr)
		}

		return fmt.Errorf("callback: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("tx commit failed", "error", err)

		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
