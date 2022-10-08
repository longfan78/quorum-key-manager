package keys

import (
	"context"

	"github.com/longfan78/quorum-key-manager/pkg/errors"
	"github.com/longfan78/quorum-key-manager/src/stores/database"

	"github.com/longfan78/quorum-key-manager/src/auth/entities"
)

func (c Connector) Restore(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring key")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey})
	if err != nil {
		return err
	}

	_, err = c.Get(ctx, id)
	if err == nil {
		return nil
	}

	_, err = c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Keys) error {
		err = dbtx.Restore(ctx, id)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support restoring, we only restore in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("key restored successfully")
	return nil
}
