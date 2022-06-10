package middleware

import (
	"bytes"
	"fmt"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/shared/utils"
	"io/ioutil"
	"regexp"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// HandlerTracing is a middleware for logging opentelemetry
func HandlerTracing(appName string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			r := c.Request()
			if utils.StringInSlice(c.Path(), []string{"/health", "/ping", "/swagger/*"}) { // exceptional don't start span
				return next(c)
			}

			// extract request header into context
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// start span
			operation := fmt.Sprintf("[API] %s %s", c.Request().Method, c.Path())
			ctx, span := otel.Tracer(appName).Start(ctx, operation)
			defer span.End()

			c.SetRequest(r.WithContext(ctx))

			// log param path
			if len(c.ParamNames()) > 0 {
				for i, param := range c.ParamNames() {
					span.SetAttributes(attribute.String(param, c.ParamValues()[i]))
				}
			}

			// log query string
			if c.QueryString() != "" {
				span.SetAttributes(attribute.String("query", c.QueryString()))
			}

			// log body
			if c.Request().Body != nil && c.Request().Header.Get("Content-Type") == "application/json" {
				body := dumpBody(c)
				span.SetAttributes(attribute.String("body", body))
			}

			span.SetAttributes(attribute.String("request_id", logger.GetRequestID(ctx)))

			return next(c)
		}
	}
}

func dumpBody(c echo.Context) string {

	// Request
	reqBody, _ := ioutil.ReadAll(c.Request().Body)

	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	ra, err := regexp.Compile("\"([a-zA-Z]*(pass|secure|token|authorization|refresh_token|access_token)[a-zA-Z]*)\"\\s?:\\s?\"([^\"]+)\"")
	if err == nil {
		return ra.ReplaceAllString(string(reqBody), "\"$1\":\"[REDACTED]\"")
	}
	return string(reqBody)

}
