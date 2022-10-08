package vaults

import (
	"context"

	"github.com/longfan78/quorum-key-manager/pkg/errors"
	auth "github.com/longfan78/quorum-key-manager/src/auth/entities"
	"github.com/longfan78/quorum-key-manager/src/entities"
	"github.com/longfan78/quorum-key-manager/src/infra/akv/client"
)

func (c *Vaults) CreateAzure(_ context.Context, name string, config *entities.AzureConfig, allowedTenants []string, _ *auth.UserInfo) error {
	logger := c.logger.With("name", name)
	logger.Debug("creating akv client")

	cli, err := client.NewClient(client.NewConfig(config))
	if err != nil {
		errMessage := "failed to instantiate AKV client"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidFormatError(errMessage)
	}

	c.createVault(name, entities.AzureVaultType, allowedTenants, cli)

	logger.Info("azure vault created successfully")
	return nil
}
