package httplog

import (
	"context"
	"go-clean-arch/infra/logger"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type HTTPLog struct {
	db *pg.DB
}

func NewHTTPLog(db *pg.DB) *HTTPLog {
	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return &HTTPLog{db}
}

func (h *HTTPLog) LogOutgoingRequest(ctx context.Context, eventName string, endpoint string, reqBody interface{}, statusCode int, res interface{}) error {
	reqID := logger.GetRequestID(ctx)
	model := LogOutgoingRequest{
		TraceID:    reqID,
		EventName:  eventName,
		Endpoint:   endpoint,
		Request:    reqBody,
		StatusCode: statusCode,
		Response:   res,
	}

	_, err := h.db.Model(&model).Insert()
	if err != nil {
		return errors.Wrap(err, "error when inserting log outgoing request")
	}
	return nil
}

func (h *HTTPLog) LogIncomingRequest(ctx context.Context, eventName, endpoint string, reqBody interface{}) error {

	reqID := logger.GetRequestID(ctx)
	model := LogIncomingRequest{
		TraceID:   reqID,
		EventName: eventName,
		Endpoint:  endpoint,
		Request:   reqBody,
	}

	_, err := h.db.Model(&model).Insert()
	if err != nil {
		return errors.Wrap(err, "error when inserting log incoming request")
	}
	return nil
}

func (h *HTTPLog) LogError(ctx context.Context, statusCode int, errorsa string, traces interface{}) error {

	reqID := logger.GetRequestID(ctx)
	model := LogError{
		TraceID:    reqID,
		StatusCode: statusCode,
		Error:      errorsa,
		Traces:     traces,
	}

	_, err := h.db.Model(&model).Insert()
	if err != nil {
		return errors.Wrap(err, "error when inserting log error")
	}
	return nil
}
