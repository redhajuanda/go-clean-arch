package middleware

import (
	"context"
	"go-clean-arch/internal/ierr"
	"go-clean-arch/shared/auth"
	"go-clean-arch/shared/response"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// VerifyJWT is a JWT middleware that verify the logged in user and set user context if verified.
// And will set user context to nil if not
func VerifyJWT(signingKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token, err := auth.VerifyTokenFromRequest(c, signingKey)
			if err != nil {
				return next(c)
			}
			c.Set("user", token)

			return next(c)
		}
	}
}

// MustLoggedIn is a JWT middleware that verify the logged in user and set user context if verified.
// And will set user context if not
func MustLoggedIn(signingKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := auth.VerifyTokenFromRequest(c, signingKey)
			if err != nil {
				return response.HTTPError(err, http.StatusUnauthorized, ierr.ErrUnauthorized.Code, ierr.ErrUnauthorized.Message)
			}

			claims := token.Claims.(jwt.MapClaims)

			var tokenType string
			if val, ok := claims["token_type"].(string); ok {
				tokenType = val
			}

			if tokenType != "access" {
				return response.ErrUnauthorized(ierr.ErrUnauthorized)
			}

			ctx := context.WithValue(c.Request().Context(), auth.ContextKeyUser, token)
			r := c.Request().WithContext(ctx)
			c.SetRequest(r)

			if ctx.Value(auth.ContextKeyUser) != nil {
				return next(c)
			}

			// else return unauthorized
			return response.ErrUnauthorized(ierr.ErrUnauthorized)

		}
	}
}
