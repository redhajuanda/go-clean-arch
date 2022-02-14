package middleware

import (
	"go-clean-arch/infra/logger"

	"github.com/labstack/echo/v4"
)

// RequestIDContext sets the context with request_id from the request if it exists, otherwise create the new request_id
func RequestIDContext() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()
			r := c.Request()

			newContext := logger.WithRequest(ctx, r)
			r = r.WithContext(newContext)

			c.SetRequest(r)
			return next(c)
		}
	}
}
