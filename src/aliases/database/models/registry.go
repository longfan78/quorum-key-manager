package models

import (
	"time"

	"github.com/longfan78/quorum-key-manager/src/entities"
)

type Registry struct {
	tableName struct{} `pg:"registries"` // nolint:unused,structcheck // reason

	Name           string    `pg:",pk"`
	Aliases        []Alias   `pg:"rel:has-many"`
	AllowedTenants []string  `pg:",array"`
	CreatedAt      time.Time `pg:"default:now()"`
	UpdatedAt      time.Time `pg:"default:now()"`
}

func NewRegistry(registry *entities.AliasRegistry) *Registry {
	return &Registry{
		Name:           registry.Name,
		AllowedTenants: registry.AllowedTenants,
		CreatedAt:      registry.CreatedAt,
		UpdatedAt:      registry.UpdatedAt,
	}
}

func (r *Registry) ToEntity() *entities.AliasRegistry {
	aliases := []entities.Alias{}
	for _, aliasModel := range r.Aliases {
		aliases = append(aliases, *aliasModel.ToEntity())
	}

	return &entities.AliasRegistry{
		Name:           r.Name,
		AllowedTenants: r.AllowedTenants,
		Aliases:        aliases,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}
