package migration

import (
	"database/sql"
	"go-clean-arch/app"
	"go-clean-arch/configs"
	"go-clean-arch/pkg/db"
	"go-clean-arch/pkg/logger"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

const (
	MIGRATION_TYPE_UP    = "up"
	MIGRATION_TYPE_DOWN  = "down"
	MIGRATION_TYPE_FRESH = "fresh"
)

type Migration struct {
	cfg *configs.Config
	log logger.Logger
	db  *sql.DB
}

func New() *Migration {
	cfg := configs.LoadDefault()
	log := logger.New(cfg.Server.NAME, app.Version)
	logger.SetFormatter(&logrus.JSONFormatter{})
	db := db.NewPostgres(cfg)

	return &Migration{
		cfg,
		log,
		db,
	}
}

func (m *Migration) Start(migrationType string) {

	m.log.Infof("start migration %s", migrationType)
	if migrationType == MIGRATION_TYPE_FRESH {
		if m.cfg.Server.ENV.IsProd() {
			m.log.Fatalf("cannot migrate fresh in production")
			return
		}
		m.log.Info("drop schema")
		_, err := m.db.Exec("drop schema public cascade; create schema public;")
		if err != nil {
			panic(err)
		}
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./scripts/migrations/postgres",
	}

	var direction migrate.MigrationDirection
	switch migrationType {
	case MIGRATION_TYPE_UP:
		direction = migrate.Up
	case MIGRATION_TYPE_DOWN:
		direction = migrate.Down
	case MIGRATION_TYPE_FRESH:
		direction = migrate.Up
	}

	count, err := migrate.Exec(m.db, "postgres", migrations, direction)
	if err != nil {
		panic(err)
	}
	m.log.Infof("applied %d migrations", count)
}
