package src

import (
	"github.com/longfan78/quorum-key-manager/pkg/http/server"
	"github.com/longfan78/quorum-key-manager/src/infra/api-key/csv"
	"github.com/longfan78/quorum-key-manager/src/infra/jwt/jose"
	"github.com/longfan78/quorum-key-manager/src/infra/log/zap"
	manifestreader "github.com/longfan78/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/longfan78/quorum-key-manager/src/infra/postgres/client"
	tls "github.com/longfan78/quorum-key-manager/src/infra/tls/filesystem"
)

type Config struct {
	HTTP     *server.Config
	Logger   *zap.Config
	Postgres *client.Config
	OIDC     *jose.Config
	APIKey   *csv.Config
	TLS      *tls.Config
	Manifest *manifestreader.Config
}
