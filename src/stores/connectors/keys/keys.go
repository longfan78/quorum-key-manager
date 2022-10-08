package keys

import (
	"github.com/longfan78/quorum-key-manager/src/auth"
	"github.com/longfan78/quorum-key-manager/src/entities"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/stores"
	"github.com/longfan78/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store        stores.KeyStore
	db           database.Keys
	logger       log.Logger
	authorizator auth.Authorizator
}

var _ stores.KeyStore = Connector{}

func NewConnector(store stores.KeyStore, db database.Keys, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		db:           db,
		logger:       logger,
		authorizator: authorizator,
	}
}

func isSupportedAlgo(alg *entities.Algorithm) bool {
	if alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1 {
		return true
	}

	if alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Babyjubjub {
		return true
	}

	if alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Curve25519 {
		return true
	}

	return false
}
