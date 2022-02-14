package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "go-clean-arch",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

func Run() {

	rootCmd.AddCommand(restCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(cronCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
