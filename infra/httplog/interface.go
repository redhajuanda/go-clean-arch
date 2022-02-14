package httplog

import "context"

type IHTTPLog interface {
	LogOutgoingRequest(ctx context.Context, eventName string, endpoint string, reqBody interface{}, statusCode int, res interface{}) error
	LogIncomingRequest(ctx context.Context, eventName, endpoint string, reqBody interface{}) error
	LogError(ctx context.Context, statusCode int, errors string, traces interface{}) error
}
