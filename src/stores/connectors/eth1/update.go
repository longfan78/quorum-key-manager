package eth1

import (
	"context"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (c Connector) Update(ctx context.Context, addr ethcommon.Address, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr.Hex())
	logger.Debug("updating ethereum account")

	err := c.authorizator.CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEth1Account})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, addr.Hex())
	if err != nil {
		return nil, err
	}
	acc.Tags = attr.Tags

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETH1Accounts) error {
		acc, err = dbtx.Update(ctx, acc)
		if err != nil {
			return err
		}

		_, err = c.store.Update(ctx, addr.Hex(), attr)
		if err != nil && !errors.IsNotSupportedError(err) {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("ethereum account updated successfully")
	return acc, nil
}
