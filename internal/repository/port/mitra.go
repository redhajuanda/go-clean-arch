package port

import (
	"context"
	"go-clean-arch/internal/domain"
)

// UserRepository encapsulates the logic to access users from the data source.
type UserRepository interface {
	// GetByID returns the user with the specified user ID.
	GetByID(ctx context.Context, userID string) (domain.User, error)
	// GetByUsername returns the user with the specified username.
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	// IsUserExist checks wether user exists
	IsUserExist(ctx context.Context, userID string) (bool, error)
	// IsUserExistByUsername checks wether user exists by username
	IsUserExistByUsername(ctx context.Context, username string) (exist bool, err error)
	// Update updates the user with given ID in the storage.
	Update(ctx context.Context, userID string, user domain.User) error
}
