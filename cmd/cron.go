package cmd

import (
	"go-clean-arch/app/cron"
	"log"

	"github.com/spf13/cobra"
)

const (
	CRON_TYPE_TRANSACTION = "transaction"
	CRON_TYPE_RECONCILE   = "reconcile"
	CRON_TYPE_CLEANUP     = "cleanup"
)

var cronCmd = &cobra.Command{
	Use: "cron",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

var cronTransactionCmd = &cobra.Command{
	Use: "transaction",
	Run: func(_ *cobra.Command, _ []string) {
		startCron(CRON_TYPE_TRANSACTION)
	},
}

var cronReconcileCmd = &cobra.Command{
	Use: "reconcile",
	Run: func(_ *cobra.Command, _ []string) {
		startCron(CRON_TYPE_RECONCILE)
	},
}

var cronCleanUpCmd = &cobra.Command{
	Use: CRON_TYPE_CLEANUP,
	Run: func(_ *cobra.Command, _ []string) {
		startCron(CRON_TYPE_CLEANUP)
	},
}

func startCron(cronType string) {
	c := cron.New()
	c.Start(cronType)
}
