package domain

import "time"

// User represents a user domain.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	FullName     *string   `json:"full_name"`
	Role         string    `json:"role"`
	RefreshToken *string   `json:"-"`
	IsActive     bool      `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetUsername returns the user username.
func (u User) GetUsername() string {
	return u.Username
}

// GetPassword return user password
func (u User) GetPassword() string {
	return u.Password
}

// GetType returns user type
func (u User) GetType() string {
	return "mitra"
}
