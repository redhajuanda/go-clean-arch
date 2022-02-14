package postgres

import (
	"context"
	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/ierr"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// UserRepository encapsulates the logic to access users from the data source.
type UserRepository struct {
	db DBI
}

// NewUserRepository creates a new user repository
func NewUserRepository(db DBI) *UserRepository {
	return &UserRepository{db}
}

// GetByID returns the user with the specified user ID.
func (r *UserRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {

	user := domain.User{ID: userID}
	q := r.db.Model(&user)
	err := q.WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return domain.User{}, ierr.ErrResourceNotFound
		}
		return domain.User{}, errors.Wrap(err, "cannot get user")
	}

	return user, nil
}

// GetByUsername returns the user with the specified username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {

	var user domain.User

	err := r.db.Model(&user).Where("username=?", username).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return domain.User{}, ierr.ErrResourceNotFound
		}
		return domain.User{}, errors.Wrap(err, "cannot get user")
	}

	return user, nil
}

// IsUserExist checks wether user exists
func (r *UserRepository) IsUserExist(ctx context.Context, userID string) (bool, error) {
	user := domain.User{ID: userID}
	exist, err := r.db.Model(&user).WherePK().Exists()
	if err != nil {
		return false, errors.Wrap(err, "cannot check user")
	}
	return exist, nil
}

// IsUserExistByUsername checks whether user exists by username
func (r *UserRepository) IsUserExistByUsername(ctx context.Context, username string) (bool, error) {
	user := domain.User{}
	exist, err := r.db.Model(&user).Where("username=?", username).Exists()
	if err != nil {
		return false, errors.Wrap(err, "cannot check user")
	}
	return exist, nil
}

// Update updates the user with given ID in the storage.
func (r *UserRepository) Update(ctx context.Context, userID string, user domain.User) error {
	user.ID = userID
	_, err := r.db.Model(&user).WherePK().UpdateNotZero()
	if err != nil {
		return errors.Wrap(err, "cannot update user")
	}
	return nil
}
