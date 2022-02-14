package user

import (
	"context"
	"go-clean-arch/internal/domain"
)

// IService encapsulates usecase logic for users.
type IService interface {
	// Get returns the user with the specified user ID or username.
	Get(ctx context.Context) (domain.User, error)
}
