package formatters

import (
	"fmt"

	common2 "github.com/longfan78/quorum-key-manager/pkg/common"
	"github.com/longfan78/quorum-key-manager/pkg/errors"
	"github.com/longfan78/quorum-key-manager/pkg/ethereum"
	"github.com/longfan78/quorum-key-manager/src/stores/api/types"
	"github.com/longfan78/quorum-key-manager/src/stores/entities"
	quorumtypes "github.com/longfan78/quorum/core/types"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

func FormatSignTypedDataRequest(request *types.SignTypedDataRequest) *signer.TypedData {
	typedData := &signer.TypedData{
		Types: signer.Types{
			ethereum.EIP712DomainLabel: []signer.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "version", Type: "string"},
			},
		},
		PrimaryType: request.MessageType,
		Domain: signer.TypedDataDomain{
			Name:              request.DomainSeparator.Name,
			Version:           request.DomainSeparator.Version,
			ChainId:           math.NewHexOrDecimal256(request.DomainSeparator.ChainID),
			VerifyingContract: request.DomainSeparator.VerifyingContract,
			Salt:              request.DomainSeparator.Salt,
		},
		Message: request.Message,
	}

	for typeName, requestTypes := range request.Types {
		var typesDefinition []signer.Type
		for _, typeDefRequest := range requestTypes {
			typesDefinition = append(typesDefinition, signer.Type{
				Name: typeDefRequest.Name,
				Type: typeDefRequest.Type,
			})
		}
		typedData.Types[typeName] = typesDefinition
	}

	if request.DomainSeparator.VerifyingContract != "" {
		typedData.Types[ethereum.EIP712DomainLabel] = append(typedData.Types[ethereum.EIP712DomainLabel], signer.Type{Name: "verifyingContract", Type: "address"})
	}

	if request.DomainSeparator.Salt != "" {
		typedData.Types[ethereum.EIP712DomainLabel] = append(typedData.Types[ethereum.EIP712DomainLabel], signer.Type{Name: "salt", Type: "string"})
	}

	return typedData
}

func FormatTransaction(tx *types.SignETHTransactionRequest) (*ethtypes.Transaction, error) {
	var txData ethtypes.TxData

	switch tx.TransactionType {
	case types.LegacyTxType:
		txData = &ethtypes.LegacyTx{
			Nonce:    uint64(tx.Nonce),
			GasPrice: tx.GasPrice.ToInt(),
			Gas:      uint64(tx.GasLimit),
			To:       tx.To,
			Value:    tx.Value.ToInt(),
			Data:     tx.Data,
		}
	case types.AccessListTxType:
		txData = &ethtypes.AccessListTx{
			ChainID:    tx.ChainID.ToInt(),
			Nonce:      uint64(tx.Nonce),
			GasPrice:   tx.GasPrice.ToInt(),
			Gas:        uint64(tx.GasLimit),
			To:         tx.To,
			Value:      tx.Value.ToInt(),
			Data:       tx.Data,
			AccessList: tx.AccessList,
		}
	case "", types.DynamicFeeTxType:
		if tx.GasFeeCap == nil {
			return nil, errors.InvalidFormatError(fmt.Sprintf("maxFeePerGas cannot be empty for a %s transaction", types.DynamicFeeTxType))
		}

		if tx.GasTipCap == nil {
			return nil, errors.InvalidFormatError(fmt.Sprintf("maxPriorityFeePerGas cannot be empty for a %s transaction", types.DynamicFeeTxType))
		}

		txData = &ethtypes.DynamicFeeTx{
			ChainID:    tx.ChainID.ToInt(),
			Nonce:      uint64(tx.Nonce),
			GasTipCap:  tx.GasTipCap.ToInt(),
			GasFeeCap:  tx.GasFeeCap.ToInt(),
			Gas:        uint64(tx.GasLimit),
			To:         tx.To,
			Value:      tx.Value.ToInt(),
			Data:       tx.Data,
			AccessList: tx.AccessList,
		}
	default:
		return nil, errors.InvalidFormatError(fmt.Sprintf("invalid transaction type, must be %s, %s or %s", types.LegacyTxType, types.AccessListTxType, types.DynamicFeeTxType))
	}

	return ethtypes.NewTx(txData), nil
}

func FormatPrivateTransaction(tx *types.SignQuorumPrivateTransactionRequest) *quorumtypes.Transaction {
	if tx.To == nil {
		return quorumtypes.NewContractCreation(uint64(tx.Nonce), tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
	}
	return quorumtypes.NewTransaction(uint64(tx.Nonce), *tx.To, tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
}

func FormatEEATransaction(tx *types.SignEEATransactionRequest) (*ethtypes.Transaction, *ethereum.PrivateArgs) {
	privateArgs := &ethereum.PrivateArgs{
		PrivateFrom: &tx.PrivateFrom,
		PrivateType: common2.ToPtr(ethereum.PrivateTypeRestricted).(*ethereum.PrivateType),
	}

	if tx.PrivacyGroupID != "" {
		privateArgs.PrivacyGroupID = &tx.PrivacyGroupID
	} else if len(tx.PrivateFor) > 0 {
		privateArgs.PrivateFor = &tx.PrivateFor
	}

	txData := &ethtypes.LegacyTx{
		Nonce:    uint64(tx.Nonce),
		GasPrice: tx.GasPrice.ToInt(),
		Gas:      uint64(tx.GasLimit),
		To:       tx.To,
		Value:    tx.Value.ToInt(),
		Data:     tx.Data,
	}

	return ethtypes.NewTx(txData), privateArgs
}

func FormatEthAccResponse(ethAcc *entities.ETHAccount) *types.EthAccountResponse {
	resp := &types.EthAccountResponse{
		KeyID:               ethAcc.KeyID,
		Address:             ethAcc.Address,
		PublicKey:           ethAcc.PublicKey,
		CompressedPublicKey: ethAcc.CompressedPublicKey,
		Tags:                ethAcc.Tags,
		CreatedAt:           ethAcc.Metadata.CreatedAt,
		UpdatedAt:           ethAcc.Metadata.UpdatedAt,
		Disabled:            ethAcc.Metadata.Disabled,
	}

	if !ethAcc.Metadata.DeletedAt.IsZero() {
		resp.DeletedAt = &ethAcc.Metadata.DeletedAt
	}

	return resp
}
