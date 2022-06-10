package cleanup

import (
	"context"
	"fmt"
	"go-clean-arch/configs"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/pkg/otel"
	"go-clean-arch/shared/utils"
	"sync"

	"github.com/go-co-op/gocron"
	"go.opentelemetry.io/otel/attribute"
)

// RegisterScheduler registers a new scheduler
func RegisterScheduler(cfg *configs.Config, log logger.Logger, service IService, cron *gocron.Scheduler, wg *sync.WaitGroup) {
	scheduler := scheduler{cfg, log, service}

	_, err := cron.Cron(cfg.Scheduler.CleanUpPattern).Do(utils.CronTaskWrapper, wg, scheduler.CleanUpHTTPLog)
	if err != nil {
		log.Fatal(err)
	}
}

type scheduler struct {
	cfg     *configs.Config
	log     logger.Logger
	service IService
}

func (s *scheduler) CleanUpHTTPLog() {

	ctx, span := otel.Start(context.Background())
	defer span.End()

	s.log.Info("execute scheduler cleanup httplog")
	err := s.service.CleanUpHTTPLog(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error_message", err.Error()))
		span.SetAttributes(attribute.String("stack_trace", fmt.Sprintf("%+v", err)))
		s.log.Error(err)
	}
}
