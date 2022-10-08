package types

import "github.com/longfan78/quorum-key-manager/src/auth/entities"

type CreateRoleRequest struct {
	Permissions []entities.Permission `json:"permissions" yaml:"permissions" validate:"required" example:"*:*"`
}
