package postgres

import (
	"context"

	"github.com/go-pg/pg/v10/orm"
)

// DBI is a DB interface implemented by *DB and *Tx.
type DBI interface {
	Model(model ...interface{}) *orm.Query
	ModelContext(c context.Context, model ...interface{}) *orm.Query
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error)
}
