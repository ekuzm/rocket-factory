package transaction

import (
	"context"
	"fmt"
)

func (m *manager) Wrap(ctx context.Context, callback func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	inject(ctx, tx)
	if err = callback(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("rollback transaction: %w", err)
		}

		return fmt.Errorf("callback: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
