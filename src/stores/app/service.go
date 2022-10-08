package app

import (
	"github.com/longfan78/quorum-key-manager/src/auth"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/infra/postgres"
	"github.com/longfan78/quorum-key-manager/src/stores/api/http"
	"github.com/longfan78/quorum-key-manager/src/stores/connectors/stores"
	db "github.com/longfan78/quorum-key-manager/src/stores/database/postgres"
	"github.com/longfan78/quorum-key-manager/src/vaults"
	"github.com/gorilla/mux"
)

func RegisterService(router *mux.Router, logger log.Logger, postgresClient postgres.Client, roles auth.Roles, vaultsService vaults.Vaults) *stores.Connector {
	// Data layer
	storesDB := db.New(logger, postgresClient)

	// Business layer
	storesService := stores.NewConnector(roles, storesDB, vaultsService, logger)

	// Service layer
	http.NewStoresHandler(storesService).Register(router)

	return storesService
}
