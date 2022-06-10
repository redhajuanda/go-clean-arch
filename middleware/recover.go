package middleware

import (
	"fmt"
	"go-clean-arch/pkg/logger"
	"runtime"

	"github.com/labstack/echo/v4"
)

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func Recover(log logger.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, !false)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])

					log.With(c.Request().Context()).Error(msg)
					c.Error(err)
				}

			}()
			return next(c)
		}
	}
}
