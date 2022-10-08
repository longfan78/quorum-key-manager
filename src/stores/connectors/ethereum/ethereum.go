package eth

import (
	"github.com/longfan78/quorum-key-manager/src/auth"
	"github.com/longfan78/quorum-key-manager/src/entities"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/stores"
	"github.com/longfan78/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store        stores.KeyStore
	logger       log.Logger
	db           database.ETHAccounts
	authorizator auth.Authorizator
}

var _ stores.EthStore = Connector{}

var ethAlgo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

func NewConnector(store stores.KeyStore, db database.ETHAccounts, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		logger:       logger,
		db:           db,
		authorizator: authorizator,
	}
}
