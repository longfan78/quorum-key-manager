package database

import (
	"context"

	"github.com/longfan78/quorum-key-manager/src/stores/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	ETHAccounts(storeID string) ETHAccounts
	Ping(ctx context.Context) error
	Keys(storeID string) Keys
	Secrets(storeID string) Secrets
}

type ETHAccounts interface {
	RunInTransaction(ctx context.Context, persistFunc func(dbtx ETHAccounts) error) error
	Get(ctx context.Context, addr string) (*entities.ETHAccount, error)
	GetDeleted(ctx context.Context, addr string) (*entities.ETHAccount, error)
	GetAll(ctx context.Context) ([]*entities.ETHAccount, error)
	GetAllDeleted(ctx context.Context) ([]*entities.ETHAccount, error)
	SearchAddresses(ctx context.Context, isDeleted bool, limit, offset uint64) ([]string, error)
	Add(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error)
	Update(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error)
	Delete(ctx context.Context, addr string) error
	Restore(ctx context.Context, addr string) error
	Purge(ctx context.Context, addr string) error
}

type Keys interface {
	RunInTransaction(ctx context.Context, persistFunc func(dbtx Keys) error) error
	Get(ctx context.Context, id string) (*entities.Key, error)
	GetDeleted(ctx context.Context, id string) (*entities.Key, error)
	GetAll(ctx context.Context) ([]*entities.Key, error)
	GetAllDeleted(ctx context.Context) ([]*entities.Key, error)
	SearchIDs(ctx context.Context, isDeleted bool, limit, offset uint64) ([]string, error)
	Add(ctx context.Context, key *entities.Key) (*entities.Key, error)
	Update(ctx context.Context, key *entities.Key) (*entities.Key, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
}

type Secrets interface {
	RunInTransaction(ctx context.Context, persistFunc func(dbtx Secrets) error) error
	Get(ctx context.Context, id, version string) (*entities.Secret, error)
	GetLatestVersion(ctx context.Context, id string, isDeleted bool) (string, error)
	ListVersions(ctx context.Context, id string, isDeleted bool) ([]string, error)
	SearchIDs(ctx context.Context, isDeleted bool, limit, offset uint64) ([]string, error)
	GetDeleted(ctx context.Context, id string) (*entities.Secret, error)
	GetAll(ctx context.Context) ([]*entities.Secret, error)
	GetAllDeleted(ctx context.Context) ([]*entities.Secret, error)
	Add(ctx context.Context, secret *entities.Secret) (*entities.Secret, error)
	Update(ctx context.Context, secret *entities.Secret) (*entities.Secret, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
}
