package models

import (
	"time"

	entities2 "github.com/longfan78/quorum-key-manager/src/entities"

	"github.com/longfan78/quorum-key-manager/src/stores/entities"
)

type Key struct {
	tableName struct{} `pg:"keys"` // nolint:unused,structcheck // reason

	ID               string `pg:",pk"`
	StoreID          string `pg:",pk"`
	PublicKey        []byte
	SigningAlgorithm string
	EllipticCurve    string
	Tags             map[string]string
	Annotations      *entities.Annotation
	Disabled         bool
	CreatedAt        time.Time `pg:"default:now()"`
	UpdatedAt        time.Time `pg:"default:now()"`
	DeletedAt        time.Time `pg:",soft_delete"`
}

func NewKey(key *entities.Key) *Key {
	return &Key{
		ID:               key.ID,
		PublicKey:        key.PublicKey,
		SigningAlgorithm: string(key.Algo.Type),
		EllipticCurve:    string(key.Algo.EllipticCurve),
		Tags:             key.Tags,
		Annotations:      key.Annotations,
		Disabled:         key.Metadata.Disabled,
		CreatedAt:        key.Metadata.CreatedAt,
		UpdatedAt:        key.Metadata.UpdatedAt,
		DeletedAt:        key.Metadata.DeletedAt,
	}
}

func (k *Key) ToEntity() *entities.Key {
	return &entities.Key{
		ID:        k.ID,
		PublicKey: k.PublicKey,
		Algo: &entities2.Algorithm{
			Type:          entities2.KeyType(k.SigningAlgorithm),
			EllipticCurve: entities2.Curve(k.EllipticCurve),
		},
		Tags:        k.Tags,
		Annotations: k.Annotations,
		Metadata: &entities.Metadata{
			Disabled:  k.Disabled,
			CreatedAt: k.CreatedAt,
			UpdatedAt: k.UpdatedAt,
			DeletedAt: k.DeletedAt,
		},
	}
}
