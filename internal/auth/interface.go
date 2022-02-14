package auth

import (
	"context"
)

// IService encapsulates the authentication logic.
type IService interface {
	// // authenticate authenticates a user using phone number and pin.
	// // It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, req RequestLogin) (LoginResponse, error)
	// RefreshToken refresh the access token
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginResponse, error)
}

// Identity represents an authenticated user iddomain.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetUsername returns the username.
	GetUsername() string
	// GetPassword returns password
	GetPassword() string
	// GetType returns user type
	GetType() string
}
