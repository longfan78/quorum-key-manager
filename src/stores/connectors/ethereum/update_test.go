package eth

import (
	"context"
	"fmt"
	"testing"

	authtypes "github.com/longfan78/quorum-key-manager/src/auth/entities"
	mock3 "github.com/longfan78/quorum-key-manager/src/auth/mock"

	"github.com/longfan78/quorum-key-manager/pkg/errors"
	"github.com/longfan78/quorum-key-manager/src/infra/log/testutils"
	"github.com/longfan78/quorum-key-manager/src/stores/database"
	mock2 "github.com/longfan78/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/longfan78/quorum-key-manager/src/stores/entities/testutils"
	"github.com/longfan78/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	acc := testutils2.FakeETHAccount()
	attributes := testutils2.FakeAttributes()
	expectedErr := fmt.Errorf("my error")
	acc.Tags = attributes.Tags

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.ETHAccounts) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should update ethAccount successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)
		store.EXPECT().Update(gomock.Any(), acc.KeyID, attributes).Return(key, nil)

		rAcc, err := connector.Update(ctx, acc.Address, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rAcc, acc)
	})

	t.Run("should update key successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)
		store.EXPECT().Update(gomock.Any(), acc.KeyID, attributes).Return(nil, rErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if key is not found", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if db fail to update", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Update(gomock.Any(), acc).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if store fail to update", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)
		store.EXPECT().Update(gomock.Any(), acc.KeyID, attributes).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
