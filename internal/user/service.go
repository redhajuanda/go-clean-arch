package user

import (
	"context"
	"go-clean-arch/config"
	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/repository/port"
	"go-clean-arch/shared/auth"
)

// Service encapsulates the user logic.
type Service struct {
	cfg         *config.Config
	repoRegitry port.RepositoryRegistry
}

// NewService creates and returns a new user service
func NewService(cfg *config.Config, repoRegitry port.RepositoryRegistry) Service {
	return Service{cfg, repoRegitry}
}

// Get returns the user with the specified user ID or username.
func (s Service) Get(ctx context.Context) (domain.User, error) {
	mitra := auth.GetLoggedInUser(ctx)

	repoUser := s.repoRegitry.GetUserRepository()
	return repoUser.GetByID(ctx, mitra.ID)
}
