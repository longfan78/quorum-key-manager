package aliases

import (
	"github.com/longfan78/quorum-key-manager/src/aliases"
	"github.com/longfan78/quorum-key-manager/src/aliases/database"
	"github.com/longfan78/quorum-key-manager/src/auth"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
)

type Aliases struct {
	aliasDB    database.Alias
	registryDB database.Registry
	logger     log.Logger
	roles      auth.Roles
}

var _ aliases.Aliases = &Aliases{}

func New(aliasDB database.Alias, registryDB database.Registry, rolesService auth.Roles, logger log.Logger) *Aliases {
	return &Aliases{
		aliasDB:    aliasDB,
		registryDB: registryDB,
		roles:      rolesService,
		logger:     logger,
	}
}
