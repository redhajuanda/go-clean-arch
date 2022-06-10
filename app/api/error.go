package api

import (
	"context"
	"fmt"
	"go-clean-arch/configs"
	"go-clean-arch/internal/ierr"
	"go-clean-arch/pkg/httplog"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/shared/response"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// CustomHTTPErrorHandler sets error response for different type of errors and logs
func CustomHTTPErrorHandler(cfg *configs.Config, log logger.Logger, httplog httplog.IHTTPLog) echo.HTTPErrorHandler {

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

		var resp response.ErrorResponse
		if res, ok := err.(response.ErrorResponse); ok {
			resp = res
		} else {
			resp = response.ErrInternalServerError(err)
		}

		if !cfg.Server.ENV.IsLocal() {
			log = log.WithStack(resp.Internal)
		}
		var traceText interface{}
		if sterr, ok := resp.Internal.(stackTracer); ok {
			traceText = fmt.Sprintf("%+v\n", sterr.StackTrace())
			if cfg.Server.ENV.IsLocal() {
				fmt.Printf("%+v\n", sterr.StackTrace())
			}
		}

		if resp.HTTPCode == 500 {
			log.With(c.Request().Context()).Error(err)
			go func(ctx context.Context) {
				err2 := httplog.LogError(ctx, 500, err.Error(), traceText)
				if err2 != nil {
					log.Error(err2)
				}
			}(c.Request().Context())
		}

		if span := trace.SpanFromContext(c.Request().Context()); span != nil {
			span.SetStatus(codes.Code(resp.HTTPCode), resp.Internal.Error())
			span.SetAttributes(attribute.Int("http.status_code", resp.HTTPCode))
			span.RecordError(errors.Cause(resp.Internal))
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.String("error_message", errors.Cause(resp.Internal).Error()))
			span.SetAttributes(attribute.String("stack_trace", fmt.Sprintf("%+v", resp.Internal)))
		}

		c.JSON(resp.HTTPCode, resp)
	}
}
