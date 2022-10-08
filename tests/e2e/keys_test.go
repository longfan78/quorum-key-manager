// +build e2e

package e2e

import (
	"encoding/base64"
	"fmt"
	utilstypes "github.com/longfan78/quorum-key-manager/src/utils/api/types"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/longfan78/quorum-key-manager/pkg/client"
	"github.com/longfan78/quorum-key-manager/pkg/common"
	"github.com/longfan78/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ecdsaPrivKey, _ = base64.StdEncoding.DecodeString("2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=")
var eddsaPrivKey, _ = base64.StdEncoding.DecodeString("X9Yz/5+O42+eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZGGbioek5qYuzJzTNZpTHrVjjFk7iFe3FYwfpxZyNPxtIaFB5gb9VP9IcHZewwNZly821re7RkmB8pGdjywygPH")

type keysTestSuite struct {
	suite.Suite
	err       error
	env       *Environment
	storeName string

	deleteQueue  *sync.WaitGroup
	destroyQueue *sync.WaitGroup
}

func TestKeyManagerKeys(t *testing.T) {
	s := new(keysTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	if len(s.env.cfg.SecretStores) == 0 {
		t.Error("list of secret stores cannot be empty")
		return
	}

	s.deleteQueue = &sync.WaitGroup{}
	s.destroyQueue = &sync.WaitGroup{}

	for _, storeN := range s.env.cfg.KeyStores {
		s.storeName = storeN
		suite.Run(t, s)
	}
}

func (s *keysTestSuite) RunT(name string, subtest func()) bool {
	return s.Run(fmt.Sprintf("%s(%s)", name, s.storeName), subtest)
}

func (s *keysTestSuite) TestCreate() {
	s.RunT("should create a new key successfully: Secp256k1/ECDSA", func() {
		keyID := fmt.Sprintf("my-key-ecdsa-%s", common.RandString(10))
		request := &types.CreateKeyRequest{
			Curve:            "secp256k1",
			SigningAlgorithm: "ecdsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		require.NoError(s.T(), err)
		defer s.queueToDelete(key)

		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), keyID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
	})

	s.RunT("should create a new key successfully: Babyjubjub/EDDSA", func() {
		keyID := fmt.Sprintf("my-key-eddsa-%s", common.RandString(10))
		request := &types.CreateKeyRequest{
			Curve:            "babyjubjub",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		// Ignoring not supported errors
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		require.NoError(s.T(), err)
		defer s.queueToDelete(key)

		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), keyID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
	})

	s.RunT("should parse errors successfully", func() {
		keyID := "my-key"
		request := &types.CreateKeyRequest{
			Curve:            "babyjubjub",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.CreateKey(s.env.ctx, "inexistentStoreName", keyID, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.RunT("should fail with bad request if curve is not supported", func() {
		keyID := "my-key"
		request := &types.CreateKeyRequest{
			Curve:            "invalidCurve",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})

	s.RunT("should fail with bad request if signing algorithm is not supported", func() {
		keyID := "my-key"
		request := &types.CreateKeyRequest{
			Curve:            "secp256k1",
			SigningAlgorithm: "invalidSigningAlgorithm",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestImport() {
	s.RunT("should create a new key successfully: Secp256k1/ECDSA", func() {
		keyID := fmt.Sprintf("my-key-import-ecdsa-%s", common.RandString(10))
		request := &types.ImportKeyRequest{
			Curve:            "secp256k1",
			PrivateKey:       ecdsaPrivKey,
			SigningAlgorithm: "ecdsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
		// Ignoring not supported errors
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		require.NoError(s.T(), err)
		defer s.queueToDelete(key)

		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R/EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a+S2RklJRFjjoLf+LI=", key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), keyID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
	})

	s.RunT("should create a new key successfully: Babyjubjub/EDDSA", func() {
		keyID := fmt.Sprintf("my-key-eddsa-%s", common.RandString(10))
		request := &types.ImportKeyRequest{
			Curve:            "babyjubjub",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
		// Ignoring not supported errors
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "X9Yz/5+O42+eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), keyID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
	})

	s.RunT("should fail with bad request if curve is not supported", func() {
		keyID := "my-key-import"
		request := &types.ImportKeyRequest{
			Curve:            "invalidCurve",
			SigningAlgorithm: "eddsa",
			PrivateKey:       ecdsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})

	s.RunT("should fail with bad request if signing algorithm is not supported", func() {
		keyID := "my-key-import"
		request := &types.ImportKeyRequest{
			Curve:            "secp256k1",
			SigningAlgorithm: "invalidSigningAlgorithm",
			PrivateKey:       ecdsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestGetKey() {
	keyID := fmt.Sprintf("my-get-key-%s", common.RandString(10))
	request := &types.ImportKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
	// Ignoring not supported errors
	if err != nil {
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
		return
	}
	defer s.queueToDelete(key)

	s.RunT("should get a key successfully", func() {
		keyRetrieved, err := s.env.client.GetKey(s.env.ctx, s.storeName, key.ID)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R/EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a+S2RklJRFjjoLf+LI=", keyRetrieved.PublicKey)
		assert.Equal(s.T(), keyID, keyRetrieved.ID)
		assert.Equal(s.T(), request.Tags, keyRetrieved.Tags)
		assert.False(s.T(), keyRetrieved.Disabled)
		assert.NotEmpty(s.T(), keyRetrieved.CreatedAt)
		assert.NotEmpty(s.T(), keyRetrieved.UpdatedAt)
	})
}

func (s *keysTestSuite) TestDeleteKey() {
	keyID := fmt.Sprintf("my-get-key-%s", common.RandString(10))
	request := &types.CreateKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
	require.NoError(s.T(), err)
	defer s.queueToDestroy(key)

	s.RunT("should delete a key successfully", func() {
		err := s.env.client.DeleteKey(s.env.ctx, s.storeName, key.ID)
		assert.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.DeleteKey(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestGetDeletedKey() {
	request := &types.CreateKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	s.RunT("should get deleted key successfully", func() {
		keyID := fmt.Sprintf("my-get-key-%s", common.RandString(10))
		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		require.NoError(s.T(), err)

		err = s.env.client.DeleteKey(s.env.ctx, s.storeName, key.ID)
		assert.NoError(s.T(), err)
		defer s.queueToDestroy(key)
		keyRetrieved, err := s.env.client.GetDeletedKey(s.env.ctx, s.storeName, key.ID)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), keyID, keyRetrieved.ID)
	})

	s.RunT("should parse errors successfully", func() {
		_, err := s.env.client.GetDeletedKey(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestRestoreKey() {
	keyID := fmt.Sprintf("my-restore-key-%s", common.RandString(10))
	request := &types.CreateKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
	require.NoError(s.T(), err)

	err = s.env.client.DeleteKey(s.env.ctx, s.storeName, key.ID)
	assert.NoError(s.T(), err)
	defer s.queueToDelete(key)

	s.RunT("should restore deleted key successfully", func() {
		errMsg := fmt.Sprintf("failed to restore key {ID: %s}", key.ID)
		err := retryOn(func() error {
			return s.env.client.RestoreKey(s.env.ctx, s.storeName, key.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)

		require.NoError(s.T(), err)

		_, err = s.env.client.GetKey(s.env.ctx, s.storeName, key.ID)
		// We should retry on status conflict for AKV
		errMsg = fmt.Sprintf("failed to get key. {ID: %s}", key.ID)
		err = retryOn(func() error {
			_, derr := s.env.client.GetKey(s.env.ctx, s.storeName, key.ID)
			return derr
		}, s.T().Logf, errMsg, http.StatusNotFound, MaxRetries)
		require.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.RestoreKey(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestDestroyKey() {
	request := &types.CreateKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	s.RunT("should destroy deleted key successfully", func() {
		keyID := fmt.Sprintf("my-restore-key-%s", common.RandString(10))
		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		require.NoError(s.T(), err)

		err = s.env.client.DeleteKey(s.env.ctx, s.storeName, key.ID)
		assert.NoError(s.T(), err)

		errMsg := fmt.Sprintf("failed to destroy key {ID: %s}", key.ID)
		err = retryOn(func() error {
			return s.env.client.DestroyKey(s.env.ctx, s.storeName, key.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)

		require.NoError(s.T(), err)

		_, err = s.env.client.GetDeletedKey(s.env.ctx, s.storeName, key.ID)
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.DestroyKey(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestListKeys() {
	keyID := fmt.Sprintf("my-key-list-%s", common.RandString(10))
	request := &types.ImportKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
	// Ignoring not supported errors
	if err != nil {
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
		return
	}
	require.NoError(s.T(), err)
	defer s.queueToDelete(key)

	s.RunT("should get all key ids successfully", func() {
		ids, err := s.env.client.ListKeys(s.env.ctx, s.storeName, 99999, 0)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, key.ID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.env.client.ListKeys(s.env.ctx, "inexistentStoreName", 0, 0)
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestListDeletedKeys() {
	keyID := fmt.Sprintf("my-deleted-key-list-%s", common.RandString(10))
	request := &types.ImportKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.env.client.ImportKey(s.env.ctx, s.storeName, keyID, request)
	// Ignoring not supported errors
	if err != nil {
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
		return
	}
	require.NoError(s.T(), err)

	err = s.env.client.DeleteKey(s.env.ctx, s.storeName, key.ID)
	assert.NoError(s.T(), err)
	defer s.queueToDestroy(key)

	s.RunT("should get all deleted key ids successfully", func() {
		ids, err := s.env.client.ListDeletedKeys(s.env.ctx, s.storeName, 99999, 0)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, key.ID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.env.client.ListKeys(s.env.ctx, "inexistentStoreName", 0, 0)
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestSignVerify() {
	data := []byte("my data to sign")
	hashedPayload := crypto.Keccak256(data)

	s.RunT("should sign a new payload successfully: Secp256k1/ECDSA", func() {
		keyID := fmt.Sprintf("my-key-sign-ecdsa-%s", common.RandString(10))
		request := &types.CreateKeyRequest{
			Curve:            "secp256k1",
			SigningAlgorithm: "ecdsa",
		}

		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		// Ignoring not supported errors
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		defer s.queueToDelete(key)

		requestSign := &types.SignBase64PayloadRequest{
			Data: hashedPayload,
		}
		signature, err := s.env.client.SignKey(s.env.ctx, s.storeName, key.ID, requestSign)
		require.NoError(s.T(), err)

		sigB, err := base64.StdEncoding.DecodeString(signature)
		require.NoError(s.T(), err)
		pubKeyB, err := base64.StdEncoding.DecodeString(key.PublicKey)
		require.NoError(s.T(), err)

		verifyRequest := &utilstypes.VerifyKeySignatureRequest{
			Data:             hashedPayload,
			Signature:        sigB,
			Curve:            key.Curve,
			SigningAlgorithm: key.SigningAlgorithm,
			PublicKey:        pubKeyB,
		}
		err = s.env.client.VerifyKeySignature(s.env.ctx, verifyRequest)
		require.NoError(s.T(), err)
	})

	s.RunT("should sign and verify a new payload successfully: Babyjubjub/EDDSA", func() {
		keyID := fmt.Sprintf("my-key-sign-eddsa-%s", common.RandString(10))
		request := &types.CreateKeyRequest{
			Curve:            "babyjubjub",
			SigningAlgorithm: "eddsa",
		}
		key, err := s.env.client.CreateKey(s.env.ctx, s.storeName, keyID, request)
		// Ignoring not supported errors
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		require.NoError(s.T(), err)
		defer s.queueToDelete(key)

		requestSign := &types.SignBase64PayloadRequest{
			Data: data,
		}
		signature, err := s.env.client.SignKey(s.env.ctx, s.storeName, key.ID, requestSign)
		require.NoError(s.T(), err)

		sigB, _ := base64.StdEncoding.DecodeString(signature)
		pubKeyB, _ := base64.StdEncoding.DecodeString(key.PublicKey)
		verifyRequest := &utilstypes.VerifyKeySignatureRequest{
			Data:             data,
			Signature:        sigB,
			Curve:            key.Curve,
			SigningAlgorithm: key.SigningAlgorithm,
			PublicKey:        pubKeyB,
		}
		err = s.env.client.VerifyKeySignature(s.env.ctx, verifyRequest)
		require.NoError(s.T(), err)
	})
}

func (s *keysTestSuite) queueToDelete(keyR *types.KeyResponse) {
	s.deleteQueue.Add(1)
	go func() {
		err := s.env.client.DeleteKey(s.env.ctx, s.storeName, keyR.ID)
		if err != nil {
			s.T().Logf("failed to delete key {ID: %s}", keyR.ID)
		} else {
			s.queueToDestroy(keyR)
		}
		s.deleteQueue.Done()
	}()
}

func (s *keysTestSuite) queueToDestroy(keyR *types.KeyResponse) {
	s.destroyQueue.Add(1)
	go func() {
		errMsg := fmt.Sprintf("failed to destroy key {ID: %s}", keyR.ID)
		err := retryOn(func() error {
			return s.env.client.DestroyKey(s.env.ctx, s.storeName, keyR.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)

		if err != nil {
			s.T().Logf(errMsg)
		}
		s.destroyQueue.Done()
	}()
}
