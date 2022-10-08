package testutils

import (
	"github.com/longfan78/quorum-key-manager/src/aliases/api/types"
)

func FakeCreateRegistryRequest() *types.CreateRegistryRequest {
	return &types.CreateRegistryRequest{
		AllowedTenants: []string{"tenant_1"},
	}
}
