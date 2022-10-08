package interceptor

import (
	"context"

	"github.com/longfan78/quorum-key-manager/pkg/jsonrpc"
	"github.com/longfan78/quorum-key-manager/src/auth/api/http"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSign(ctx context.Context, from ethcommon.Address, data hexutil.Bytes) (*hexutil.Bytes, error) {
	logger := i.logger.With("from_account", from.Hex())
	logger.Debug("signing payload")

	store, err := i.stores.EthereumByAddr(ctx, from, http.UserInfoFromContext(ctx))
	if err != nil {
		return nil, err
	}

	sig, err := store.Sign(ctx, from, data)
	if err != nil {
		return nil, err
	}

	logger.Info("payload signed successfully")
	return (*hexutil.Bytes)(&sig), nil
}

func (i *Interceptor) EthSign() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSign)
	return h
}
