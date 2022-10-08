package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/longfan78/quorum-key-manager/src/auth/entities"
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

func TestDestroyKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.Keys) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should destroy key successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Purge(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), key.ID).Return(nil)

		err := connector.Destroy(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should destroy key successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Purge(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), key.ID).Return(rErr)

		err := connector.Destroy(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(expectedErr)

		err := connector.Destroy(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy key if key is not deleted", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, expectedErr)

		err := connector.Destroy(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy key if db fail to purge", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Purge(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Destroy(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy key if store fail to destroy", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Purge(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Destroy(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
