package user

import (
	"context"
	"go-clean-arch/configs"
	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/repository/port"
	"go-clean-arch/pkg/otel"
	"go-clean-arch/shared/auth"
)

// Service encapsulates the user logic.
type Service struct {
	cfg         *configs.Config
	repoRegitry port.RepositoryRegistry
}

// NewService creates and returns a new user service
func NewService(cfg *configs.Config, repoRegitry port.RepositoryRegistry) Service {
	return Service{cfg, repoRegitry}
}

// Get returns the user with the specified user ID or username.
func (s Service) Get(ctx context.Context) (domain.User, error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	mitra := auth.GetLoggedInUser(ctx)

	repoUser := s.repoRegitry.GetUserRepository()
	return repoUser.GetByID(ctx, mitra.ID)
}
