package secrets

import (
	"github.com/longfan78/quorum-key-manager/src/auth"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/stores"
	"github.com/longfan78/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store        stores.SecretStore
	logger       log.Logger
	db           database.Secrets
	authorizator auth.Authorizator
}

var _ stores.SecretStore = &Connector{}

func NewConnector(store stores.SecretStore, db database.Secrets, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		logger:       logger,
		db:           db,
		authorizator: authorizator,
	}
}
