package cron

import (
	"fmt"
	"go-clean-arch/app"
	"go-clean-arch/configs"
	"go-clean-arch/internal/cleanup"
	"go-clean-arch/internal/repository/postgres"
	"go-clean-arch/pkg/db"
	"go-clean-arch/pkg/httplog"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/pkg/otel"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

const (
	CRON_TYPE_CLEANUP = "cleanup"
)

type Cron struct {
	cfg *configs.Config
	log logger.Logger
	db  *pg.DB
}

func New() *Cron {
	cfg := configs.LoadDefault()
	log := logger.New(cfg.Server.NAME, app.Version)
	logger.SetFormatter(&logrus.JSONFormatter{})
	db := db.NewGoPG(cfg)

	return &Cron{
		cfg,
		log,
		db,
	}
}

func (c *Cron) Start(cronType string) {

	service := fmt.Sprintf("%s-cron-%s", c.cfg.Server.NAME, cronType)
	err := otel.SetTraceProvider(c.cfg.OpenTelemetry.JaegerURL, service, app.Version, c.cfg.Server.ENV.String(), c.cfg.OpenTelemetry.Sampled)
	if err != nil {
		c.log.Fatal(err)
	}

	repoRegistry := postgres.NewRepositoryRegistry(c.db)

	// new scheduler
	cron := gocron.NewScheduler(time.Local)
	wg := &sync.WaitGroup{}

	switch cronType {
	case CRON_TYPE_CLEANUP:
		cleanUpSvc := cleanup.NewService(c.cfg, c.log, repoRegistry, httplog.NewHTTPLog(c.db))
		// register scheduler
		cleanup.RegisterScheduler(c.cfg, c.log, cleanUpSvc, cron, wg)

	default:
		c.log.Fatalf("no cron type available")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	cron.StartAt(time.Now())
	cron.StartAsync()
	c.log.Infof("cron %s is running", cronType)

	<-signalChan

	c.log.Info("got signal to exit cron")
	cron.Clear()
	wg.Wait()

	c.log.Info("exiting cron gracefully")
}
