package vaults

import (
	"context"

	"github.com/longfan78/quorum-key-manager/pkg/errors"
	auth "github.com/longfan78/quorum-key-manager/src/auth/entities"
	"github.com/longfan78/quorum-key-manager/src/entities"
	"github.com/longfan78/quorum-key-manager/src/infra/aws/client"
)

func (c *Vaults) CreateAWS(_ context.Context, name string, config *entities.AWSConfig, allowedTenants []string, _ *auth.UserInfo) error {
	logger := c.logger.With("name", name)
	logger.Debug("creating aws vault client")

	cli, err := client.New(client.NewConfig(config), logger)
	if err != nil {
		errMessage := "failed to instantiate AWS client"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.createVault(name, entities.AWSVaultType, allowedTenants, cli)

	logger.Info("aws vault created successfully")
	return nil
}
