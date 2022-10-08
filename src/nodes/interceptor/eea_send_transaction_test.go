package interceptor

import (
	"context"
	"math/big"
	"testing"

	"github.com/longfan78/quorum-key-manager/src/auth/api/http"

	"github.com/longfan78/quorum-key-manager/pkg/common"
	"github.com/longfan78/quorum-key-manager/pkg/ethereum"
	mockethereum "github.com/longfan78/quorum-key-manager/pkg/ethereum/mock"
	"github.com/longfan78/quorum-key-manager/src/auth/entities"
	proxynode "github.com/longfan78/quorum-key-manager/src/nodes/node/proxy"
	mockaccounts "github.com/longfan78/quorum-key-manager/src/stores/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEEASendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores, aliases := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockEthStore(ctrl)

	userInfo := &entities.UserInfo{
		Username:    "username",
		Roles:       []string{"role1", "role2"},
		Permissions: []entities.Permission{"write:key", "read:key", "sign:key"},
	}
	session := proxynode.NewMockSession(ctrl)
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = http.WithUserInfo(ctx, userInfo)

	cller := mockethereum.NewMockCaller(ctrl)
	eeaCaller := mockethereum.NewMockEEACaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	privCaller := mockethereum.NewMockPrivCaller(ctrl)
	cller.EXPECT().EEA().Return(eeaCaller).AnyTimes()
	cller.EXPECT().Eth().Return(ethCaller).AnyTimes()
	cller.EXPECT().Priv().Return(privCaller).AnyTimes()

	session.EXPECT().EthCaller().Return(cller).AnyTimes()

	tests := []*testHandlerCase{
		{
			desc:    "Transaction with Privacy Get ID",
			handler: i,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eea_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","privacyGroupId":"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}],"id":"abcd"}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().EthereumByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				ethCaller.EXPECT().EstimateGas(gomock.Any(), gomock.Any()).Return(uint64(21000), nil)

				// Get Nonc
				privCaller.EXPECT().GetTransactionCount(gomock.Any(), expectedFrom, "kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=").Return(uint64(5), nil)

				// SignEEA
				expectedPrivateArgs := (&ethereum.PrivateArgs{PrivateType: common.ToPtr(ethereum.PrivateTypeRestricted).(*ethereum.PrivateType)}).WithPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
				accountsStore.EXPECT().SignEEA(gomock.Any(), expectedFrom, big.NewInt(1998), gomock.Any(), expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				// SendRawTransaction
				eeaCaller.EXPECT().SendRawTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27")).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)

				aliases.EXPECT().Replace(gomock.Any(), []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}, userInfo).Return([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}, nil)
				aliases.EXPECT().ReplaceSimple(gomock.Any(), "GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=", userInfo).Return("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=", nil)
				aliases.EXPECT().Parse("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=").Return("", "", false)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":"abcd"}`),
		},
		{
			desc:    "Transaction with privateFor",
			handler: i,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eea_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}],"id":"abcd"}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().EthereumByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				ethCaller.EXPECT().EstimateGas(gomock.Any(), gomock.Any()).Return(uint64(21000), nil)

				// Get Nonc
				privCaller.EXPECT().GetEeaTransactionCount(gomock.Any(), expectedFrom, "GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=", []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}).Return(uint64(5), nil)

				// Sign
				expectedPrivateArgs := (&ethereum.PrivateArgs{PrivateType: common.ToPtr(ethereum.PrivateTypeRestricted).(*ethereum.PrivateType)}).WithPrivateFrom("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=").WithPrivateFor([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="})
				accountsStore.EXPECT().SignEEA(gomock.Any(), expectedFrom, big.NewInt(1998), gomock.Any(), expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				eeaCaller.EXPECT().SendRawTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27")).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":"abcd"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
