package transaction

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
)

func (m *manager) Wrap(ctx context.Context, callback func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create transaction")

		return fmt.Errorf("create transaction: %w", err)
	}

	inject(ctx, tx)
	if err = callback(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			logger.WithFields(logrus.Fields{
				"Callback Error": err,
				"error":          rollbackErr,
			}).Error("Failed to rollback transaction")

			return fmt.Errorf("rollback transaction: %w", rollbackErr)
		}

		return fmt.Errorf("callback: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to commit transaction")

		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
