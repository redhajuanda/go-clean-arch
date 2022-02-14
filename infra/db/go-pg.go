package db

import (
	"context"
	"fmt"
	"go-clean-arch/config"
	"log"

	"github.com/go-pg/pg/v10"
)

// NewGoPG creates a new postgres connection
func NewGoPG(cfg *config.Config) *pg.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)

	opt, err := pg.ParseURL(connString)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)

	if cfg.Server.DEBUG {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	bt, _ := q.FormattedQuery()
	log.Printf(string(bt))
	return nil
}
