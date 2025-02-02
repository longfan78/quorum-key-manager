package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/longfan78/quorum-key-manager/src/auth/entities"
	mock3 "github.com/longfan78/quorum-key-manager/src/auth/mock"

	"github.com/longfan78/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/longfan78/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/longfan78/quorum-key-manager/src/stores/entities/testutils"
	"github.com/longfan78/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []byte("0x123")
	result := []byte("0x456")
	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should encrypt data successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceKey}).Return(nil)
		store.EXPECT().Encrypt(gomock.Any(), key.ID, data).Return(result, nil)

		rResult, err := connector.Encrypt(ctx, key.ID, data)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceKey}).Return(expectedErr)

		_, err := connector.Encrypt(ctx, key.ID, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to encrypt data if encrypt fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceKey}).Return(nil)
		store.EXPECT().Encrypt(gomock.Any(), key.ID, data).Return(nil, expectedErr)

		_, err := connector.Encrypt(ctx, key.ID, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
