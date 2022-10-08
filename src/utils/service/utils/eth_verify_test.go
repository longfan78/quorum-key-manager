package utils

import (
	"testing"

	"github.com/longfan78/quorum-key-manager/pkg/errors"

	"github.com/longfan78/quorum-key-manager/src/infra/log/testutils"
	testutils2 "github.com/longfan78/quorum-key-manager/src/stores/entities/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVerifyMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	acc := testutils2.FakeETHAccount()
	acc.Address = common.HexToAddress("0x185bD93d8D62AF4e7AD6c928561A3d86142e26ef")

	t.Run("should verify message successfully with recID 27", func(t *testing.T) {
		data := crypto.Keccak256([]byte("my data to sign"))
		address := common.HexToAddress("0xf9602d642310A014048b8325eeF3743214b9f36a")
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE1B")

		err := connector.VerifyMessage(address, data, ecdsaSignature)

		assert.NoError(t, err)
	})

	t.Run("should verify message successfully with recID 28", func(t *testing.T) {
		data := crypto.Keccak256([]byte("my data to sign"))
		address := common.HexToAddress("0x185bD93d8D62AF4e7AD6c928561A3d86142e26ef")
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE1C")

		err := connector.VerifyMessage(address, data, ecdsaSignature)

		assert.NoError(t, err)
	})

	t.Run("should fail to verify message if signature length is incorrect", func(t *testing.T) {
		data := crypto.Keccak256([]byte("xxxx"))
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE")

		err := connector.VerifyMessage(acc.Address, data, ecdsaSignature)

		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Equal(t, err.Error(), "IR500: signature must be exactly 65 bytes")
	})

	t.Run("should fail to verify message if recovery ID is incorrect", func(t *testing.T) {
		data := crypto.Keccak256([]byte("xxxx"))
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE08")

		err := connector.VerifyMessage(acc.Address, data, ecdsaSignature)

		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Equal(t, "IR500: invalid signature recovery id", err.Error())
	})

	t.Run("should fail to verify message if data is incorrect", func(t *testing.T) {
		data := crypto.Keccak256([]byte("xxxx"))
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE1C")

		err := connector.VerifyMessage(acc.Address, data, ecdsaSignature)

		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Equal(t, "IR500: failed to verify signature: recovered address does not match the expected one or payload is malformed", err.Error())
	})
}
