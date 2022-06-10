package cleanup

import (
	"context"
	"go-clean-arch/configs"
	"go-clean-arch/internal/repository/port"
	"go-clean-arch/pkg/httplog"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/shared/times"
)

// Service encapsulates the transaction logic.
type Service struct {
	cfg          *configs.Config
	log          logger.Logger
	repoRegistry port.RepositoryRegistry
	httplog      httplog.IHTTPLog
}

// NewService creates and returns a new transaction service
func NewService(cfg *configs.Config, log logger.Logger, repoRegistry port.RepositoryRegistry, httplog httplog.IHTTPLog) *Service {
	return &Service{cfg, log, repoRegistry, httplog}
}

func (s *Service) CleanUpHTTPLog(ctx context.Context) error {

	return s.httplog.CleanUpLog(ctx, times.Now().AddDate(0, 0, -7))
}
