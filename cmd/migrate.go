package cmd

import (
	"go-clean-arch/app/migration"
	"log"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("down")
	},
}

var migrateFreshCmd = &cobra.Command{
	Use: "fresh",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("fresh")
	},
}

func startMigrate(migrationType string) {
	m := migration.New()
	m.Start(migrationType)
}
