package transaction

import (
	"context"
	"fmt"
)

func Wrap(ctx context.Context, fn func(context.Context) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = fn(context.WithValue(ctx, ctxKey{}, tx)); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("failed to transaction rollback: %w", err)
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to transaction commit: %w", err)
	}

	return nil
}
