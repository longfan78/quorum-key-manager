package utils

import (
	"github.com/longfan78/quorum-key-manager/src/infra/log"
	"github.com/longfan78/quorum-key-manager/src/utils"
)

type Utilities struct {
	logger log.Logger
}

var _ utils.Utilities = &Utilities{}

func New(logger log.Logger) *Utilities {
	return &Utilities{
		logger: logger,
	}
}
