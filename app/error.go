package app

import (
	"context"
	"fmt"
	"go-clean-arch/config"
	"go-clean-arch/infra/httplog"
	"go-clean-arch/infra/logger"
	"go-clean-arch/internal/ierr"
	"go-clean-arch/shared/response"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// CustomHTTPErrorHandler sets error response for different type of errors and logs
func CustomHTTPErrorHandler(cfg *config.Config, log logger.Logger, httplog httplog.IHTTPLog) echo.HTTPErrorHandler {

	return func(err error, c echo.Context) {

		// var internalErr error
		if _, ok := err.(response.ErrorResponse); !ok {
			err = response.ErrorResponse{
				Success:   false,
				HTTPCode:  http.StatusInternalServerError,
				Message:   ierr.ErrInternal.Message,
				ErrorCode: ierr.ErrInternal.Code,
				Internal:  err,
			}
		}

		// Get internal error
		internalErr := errors.Cause(err.(response.ErrorResponse).Internal)

		// handles resource not found errors
		if errors.Is(internalErr, echo.ErrNotFound) {
			err = response.HTTPError(internalErr, http.StatusNotFound, ierr.ErrResourceNotFound.Code, "requested endpoint is not registered")
		}

		// Handles validation error
		if errors.As(internalErr, &validation.Errors{}) || errors.As(internalErr, &validation.ErrorObject{}) {
			err = response.HTTPError(internalErr, http.StatusBadRequest, ierr.ErrBadRequest.Code, internalErr.Error())
		}

		if resp, ok := err.(response.ErrorResponse); ok {

			var traceText interface{}
			if sterr, ok := resp.Internal.(stackTracer); ok {
				traceText = fmt.Sprintf("%+v\n", sterr.StackTrace())
				if cfg.Server.ENV == "local" {
					fmt.Printf("%+v\n", sterr.StackTrace())
				}
			}
			if !cfg.Server.ENV.IsLocal() {
				log = log.WithStack(resp.Internal)
			}

			log.With(c.Request().Context()).Error(resp.Internal)

			if resp.HTTPCode == 500 {
				go func(ctx context.Context) {
					err2 := httplog.LogError(ctx, resp.HTTPCode, resp.Internal.Error(), traceText)
					if err2 != nil {
						log.Error(err2)
					}
				}(c.Request().Context())
			}
			c.JSON(resp.HTTPCode, resp)
			return
		} else {
			log.With(c.Request().Context()).Error(err)

			go func(ctx context.Context) {
				err2 := httplog.LogError(ctx, 500, err.Error(), nil)
				if err2 != nil {
					log.Error(err2)
				}
			}(c.Request().Context())

			c.JSON(http.StatusInternalServerError, response.ErrInternalServerError(err))
		}
	}
}
