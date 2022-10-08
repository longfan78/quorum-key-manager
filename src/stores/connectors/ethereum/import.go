package eth

import (
	"context"

	"github.com/longfan78/quorum-key-manager/src/stores/database/models"

	"github.com/longfan78/quorum-key-manager/pkg/errors"

	authentities "github.com/longfan78/quorum-key-manager/src/auth/entities"

	"github.com/longfan78/quorum-key-manager/src/stores/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETHAccount, error) {
	logger := c.logger.With("id", id)
	logger.Debug("importing ethereum account")

	if len(privKey) == 0 {
		errMessage := "private key must be provided"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionWrite, Resource: authentities.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	key, err := c.store.Import(ctx, id, privKey, ethAlgo, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		key, err = c.store.Get(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, models.NewETHAccountFromKey(key, attr))
	if err != nil {
		return nil, err
	}

	logger.With("address", acc.Address, "key_id", acc.KeyID).Info("ethereum account imported successfully")
	return acc, nil
}
