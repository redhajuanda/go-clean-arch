package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "application",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

func Run() {

	// api
	rootCmd.AddCommand(apiCmd)

	// migrate
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(migrateCmd)

	//cron
	cronCmd.AddCommand(cronTransactionCmd)
	cronCmd.AddCommand(cronReconcileCmd)
	cronCmd.AddCommand(cronCleanUpCmd)
	rootCmd.AddCommand(cronCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
