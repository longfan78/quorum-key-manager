package keys

import (
	"context"

	entities2 "github.com/longfan78/quorum-key-manager/src/entities"

	"github.com/longfan78/quorum-key-manager/pkg/errors"

	authentities "github.com/longfan78/quorum-key-manager/src/auth/entities"

	"github.com/longfan78/quorum-key-manager/src/stores/entities"
)

func (c Connector) Create(ctx context.Context, id string, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id, "algorithm", alg.Type, "curve", alg.EllipticCurve)
	logger.Debug("creating key")

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionWrite, Resource: authentities.ResourceKey})
	if err != nil {
		return nil, err
	}

	if !isSupportedAlgo(alg) {
		errMessage := "invalid or not supported elliptic curve and signing algorithm combination for creation"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	key, err := c.store.Create(ctx, id, alg, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		key, err = c.store.Get(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	key, err = c.db.Add(ctx, key)
	if err != nil {
		return nil, err
	}

	logger.Info("key created successfully")
	return key, nil
}
