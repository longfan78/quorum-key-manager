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

func TestRestoreKey(t *testing.T) {
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

	t.Run("should restore key successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)

		err := connector.Restore(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should be idempotent when key already exists", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, nil)

		err := connector.Restore(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should restore key successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(rErr)

		err := connector.Restore(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail if key not deleted yet", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(nil, expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore key if db fail to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore key if store fail to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceKey}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
