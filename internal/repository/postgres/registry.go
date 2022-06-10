package postgres

import (
	"context"
	"go-clean-arch/internal/repository/port"
	"go-clean-arch/pkg/otel"

	"github.com/go-pg/pg/v10"
)

type RepositoryRegistry struct {
	db         *pg.DB
	dbExecutor DBI
}

func NewRepositoryRegistry(db *pg.DB) port.RepositoryRegistry {
	repo := &RepositoryRegistry{
		db: db,
	}
	return repo
}

func (r *RepositoryRegistry) DoInTransaction(ctx context.Context, txFunc port.InTransaction) (out interface{}, err error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	var tx *pg.Tx
	registry := r
	if r.dbExecutor == nil {
		tx, err = r.db.BeginContext(ctx)
		if err != nil {
			return
		}

		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				rErr := tx.Rollback() // err is non-nil; don't change it
				if rErr != nil {
					err = rErr
				}
			} else {
				err = tx.Commit() // err is nil; if Commit returns error update err
			}
		}()

		registry = &RepositoryRegistry{
			db:         r.db,
			dbExecutor: tx,
		}
	}

	out, err = txFunc(ctx, registry)

	return
}

func (r *RepositoryRegistry) GetUserRepository() port.UserRepository {
	if r.dbExecutor != nil {
		return NewUserRepository(r.dbExecutor)
	}
	return NewUserRepository(r.db)
}
