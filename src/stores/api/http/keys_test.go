package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	entities2 "github.com/longfan78/quorum-key-manager/src/entities"
	"github.com/longfan78/quorum-key-manager/src/stores/api/formatters"

	authapi "github.com/longfan78/quorum-key-manager/src/auth/api/http"

	"github.com/longfan78/quorum-key-manager/pkg/errors"
	authentities "github.com/longfan78/quorum-key-manager/src/auth/entities"
	http2 "github.com/longfan78/quorum-key-manager/src/infra/http"
	"github.com/longfan78/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/longfan78/quorum-key-manager/src/stores/entities"
	testutils2 "github.com/longfan78/quorum-key-manager/src/stores/entities/testutils"
	"github.com/longfan78/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	keyStoreName            = "KeyStore"
	invalidSigningAlgorithm = "invalidSigningAlgorithm"
	invalidCurve            = "invalidCurve"
	keyID                   = "my-key"
)

var keyUserInfo = &authentities.UserInfo{
	Username:    "username",
	Roles:       []string{"role1", "role2"},
	Permissions: []authentities.Permission{"write:key", "read:key", "sign:key"},
}

type keysHandlerTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	stores   *mock.MockStores
	keyStore *mock.MockKeyStore
	router   *mux.Router
	ctx      context.Context
}

func TestKeysHandler(t *testing.T) {
	s := new(keysHandlerTestSuite)
	suite.Run(t, s)
}

func (s *keysHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.stores = mock.NewMockStores(s.ctrl)
	s.keyStore = mock.NewMockKeyStore(s.ctrl)

	s.stores.EXPECT().Key(gomock.Any(), keyStoreName, keyUserInfo).Return(s.keyStore, nil).AnyTimes()

	s.router = mux.NewRouter()
	s.ctx = authapi.WithUserInfo(context.Background(), keyUserInfo)
	NewStoresHandler(s.stores).Register(s.router)
}

func (s *keysHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *keysHandlerTestSuite) TestCreate() {
	s.Run("should execute request successfully", func() {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Create(
			gomock.Any(),
			keyID,
			&entities2.Algorithm{
				Type:          entities2.KeyType(createKeyRequest.SigningAlgorithm),
				EllipticCurve: entities2.Curve(createKeyRequest.Curve),
			},
			&entities.Attributes{
				Tags: createKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with no body successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID, nil).WithContext(s.ctx)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
			&entities2.Algorithm{},
			&entities.Attributes{}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should fail with 400 if signing algorithm is not supported", func() {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		createKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	s.Run("should fail with 400 if curve is not supported", func() {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		createKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.keyStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestImport() {
	s.Run("should execute request successfully", func() {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID+"/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		key := testutils2.FakeKey()
		s.keyStore.EXPECT().Import(
			gomock.Any(),
			keyID,
			importKeyRequest.PrivateKey,
			&entities2.Algorithm{
				Type:          entities2.KeyType(importKeyRequest.SigningAlgorithm),
				EllipticCurve: entities2.Curve(importKeyRequest.Curve),
			},
			&entities.Attributes{
				Tags: importKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should fail with 400 if signing algorithm is not supported", func() {
		importKeyRequest := testutils.FakeImportKeyRequest()
		importKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID+"/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	s.Run("should fail with 400 if curve is not supported", func() {
		importKeyRequest := testutils.FakeImportKeyRequest()
		importKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID+"/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/"+keyID+"/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.keyStore.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestSign() {
	s.Run("should execute request successfully", func() {
		signPayloadRequest := testutils.FakeSignBase64PayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/KeyStore/keys/%s/sign", keyID), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signature := []byte("signature")
		s.keyStore.EXPECT().Sign(gomock.Any(), keyID, signPayloadRequest.Data, gomock.Any()).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), base64.URLEncoding.EncodeToString(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signPayloadRequest := testutils.FakeSignBase64PayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/KeyStore/keys/%s/sign", keyID), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.keyStore.EXPECT().Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestGet() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil).WithContext(s.ctx)

		key := testutils2.FakeKey()
		s.keyStore.EXPECT().Get(gomock.Any(), keyID).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil).WithContext(s.ctx)

		s.keyStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/KeyStore/keys", nil).WithContext(s.ctx)

		ids := []string{"key1", "key2"}
		s.keyStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: ids,
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with limit and offset successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/KeyStore/keys?limit=5&page=2", nil).WithContext(s.ctx)

		ids := []string{"key1", "key2"}
		s.keyStore.EXPECT().List(gomock.Any(), uint64(5), uint64(10)).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: ids,
			Paging: http2.PagePagingResponse{
				Previous: "example.com?limit=5&page=1",
			},
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/KeyStore/keys", nil).WithContext(s.ctx)

		s.keyStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestDelete() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil).WithContext(s.ctx)

		s.keyStore.EXPECT().Delete(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should execute request successfully with version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil).WithContext(s.ctx)
		httpRequest = httpRequest.WithContext(WithStoreName(httpRequest.Context(), keyStoreName))

		s.keyStore.EXPECT().Delete(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil).WithContext(s.ctx)

		s.keyStore.EXPECT().Delete(gomock.Any(), keyID).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestDestroy() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s/destroy", keyID), nil).WithContext(s.ctx)

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should execute request successfully with version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s/destroy", keyID), nil).WithContext(s.ctx)
		httpRequest = httpRequest.WithContext(WithStoreName(httpRequest.Context(), keyStoreName))

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s/destroy", keyID), nil).WithContext(s.ctx)

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}
