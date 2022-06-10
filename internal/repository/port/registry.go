package port

import (
	"context"
)

type InTransaction func(ctx context.Context, repoRegistry RepositoryRegistry) (interface{}, error)

type RepositoryRegistry interface {
	DoInTransaction(ctx context.Context, txFunc InTransaction) (out interface{}, err error)
	GetUserRepository() UserRepository
}
