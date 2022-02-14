package cmd

import (
	"go-clean-arch/config"
	"go-clean-arch/infra/db"
	"go-clean-arch/infra/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use: "cron",
	Run: func(_ *cobra.Command, _ []string) {

		cfg := config.LoadDefault()

		log := logger.New(cfg.Server.NAME, Version)
		logger.SetFormatter(&logrus.JSONFormatter{})

		db := db.NewGoPG(cfg)
		defer func() {
			log.Info("closing db connection")
			db.Close()
		}()

		// new scheduler
		cron := gocron.NewScheduler()
		wg := &sync.WaitGroup{}

		// register scheduler
		// payment.RegisterScheduler(cfg, log, paymentSvc, cron, wg)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		cron.Start()
		log.Info("cron is running")
		log.Info("jobs", cron.Jobs())

		<-signalChan

		cron.Clear()
		wg.Wait()

		log.Info("exiting gracefully")
	},
}
