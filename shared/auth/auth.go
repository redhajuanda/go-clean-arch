package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

type ContextUser string

const ContextKeyUser ContextUser = "user"

// GetLoggedInUser returns logged in user crm
func GetLoggedInUser(ctx context.Context) User {

	v := ctx.Value(ContextKeyUser)
	if v == nil {
		return User{}
	}
	loggedInUser := v.(*jwt.Token)

	claims := loggedInUser.Claims.(jwt.MapClaims)

	var id string
	if val, ok := claims["id"].(string); ok {
		id = val
	}

	var username string
	if val, ok := claims["username"].(string); ok {
		username = val
	}

	var role string
	if val, ok := claims["user_type"].(string); ok {
		role = val
	}

	return User{
		ID:       id,
		Username: username,
		Role:     role,
	}

}
