package keys

import (
	"context"

	"github.com/longfan78/quorum-key-manager/src/auth/entities"

	"github.com/longfan78/quorum-key-manager/pkg/errors"

	"github.com/longfan78/quorum-key-manager/src/stores/database"
)

func (c Connector) Delete(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("deleting key")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey})
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Keys) error {
		derr := dbtx.Delete(ctx, id)
		if derr != nil {
			return derr
		}

		derr = c.store.Delete(ctx, id)
		if derr != nil && !errors.IsNotSupportedError(derr) { // If the underlying store does not support deleting, we only delete in DB
			return derr
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}
