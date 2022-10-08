package e2e

import (
	"context"
	"github.com/longfan78/quorum-key-manager/pkg/client"
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/infra/log/zap"
	"github.com/longfan78/quorum-key-manager/tests"
	"net/http"
)

type Environment struct {
	ctx        context.Context
	logger     log.Logger
	httpClient *http.Client
	client     client.KeyManagerClient
	cfg        *tests.Config
}

func NewEnvironment() (*Environment, error) {
	cfg, err := tests.NewConfig()
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewLogger(zap.NewConfig(zap.InfoLevel, zap.JSONFormat))
	if err != nil {
		return nil, err
	}

	keyManagerClient := client.NewHTTPClient(
		&http.Client{Transport: NewTestHttpTransport("", cfg.ApiKey, nil)},
		&client.Config{URL: cfg.KeyManagerURL},
	)

	return &Environment{
		ctx:    context.Background(),
		logger: logger,
		client: keyManagerClient,
		cfg:    cfg,
	}, nil
}
