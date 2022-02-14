package cmd

import (
	"go-clean-arch/config"
	"go-clean-arch/infra/db"
	"log"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(_ *cobra.Command, _ []string) {
		log.Fatal("no subcommand found")
	},
}

var migrateUpCmd = &cobra.Command{
	Use: "up",

	Run: func(_ *cobra.Command, _ []string) {
		log.Println("migrate up")
		cfg := config.LoadDefault()
		migration(cfg, false)
	},
}

var migrateFreshCmd = &cobra.Command{
	Use: "fresh",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := config.LoadDefault()
		if cfg.Server.ENV.IsProd() {
			log.Fatalln("cannot migrate fresh in production environment")
		}
		log.Println("migrate fresh")
		migration(cfg, true)
	},
}

func migration(cfg *config.Config, fresh bool) {
	sql := db.NewPostgres(cfg)
	if fresh {
		log.Println("drop schema")
		_, err := sql.Exec("drop schema public cascade; create schema public;")
		if err != nil {
			panic(err)
		}
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./scripts/migrations/postgres",
	}

	count, err := migrate.Exec(sql, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	log.Println("applied", count, "migrations")
}
