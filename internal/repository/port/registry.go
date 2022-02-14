package port

import (
	"context"
)

type InTransaction func(repoRegistry RepositoryRegistry) (interface{}, error)

type RepositoryRegistry interface {
	DoInTransaction(ctx context.Context, txFunc InTransaction) (out interface{}, err error)
	GetUserRepository() UserRepository
}
