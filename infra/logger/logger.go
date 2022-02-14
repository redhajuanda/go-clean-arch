package logger

import (
	"context"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Params type, used to pass to `WithParams`.
type Params map[string]interface{}

// Logger represent common interface for logging function
type Logger interface {
	With(ctx context.Context) Logger
	WithStack(err error) Logger
	WithParam(key string, value interface{}) Logger
	WithParams(params Params) Logger
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

type logger struct {
	*logrus.Entry
}

var logStore *logger

// New returns a new wrapper log
func New(serviceName, serviceVersion string) Logger {
	logStore = &logger{logrus.New().WithFields(logrus.Fields{"service": serviceName, "version": serviceVersion})}
	return logStore
}

// SetOutput sets the logger output.
func SetOutput(output io.Writer) {
	logStore.Logger.SetOutput(output)
}

// SetFormatter sets the logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logStore.Logger.SetFormatter(formatter)
}

// SetLevel sets the logger level.
func SetLevel(level logrus.Level) {
	logStore.Logger.SetLevel(logrus.InfoLevel)
}

// With reads requestId and correlationId from context and adds to log field
func (l *logger) With(ctx context.Context) Logger {

	var le *logrus.Entry
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			le = l.WithField("request_id", id)
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			le = l.WithField("correlation_id", id)
		}
	}
	return &logger{le}

}

func (l *logger) WithStack(err error) Logger {

	stack := MarshalStack(err)
	return &logger{l.WithField("stack", stack)}
}

func (l *logger) WithParam(key string, value interface{}) Logger {

	return &logger{l.WithField(key, value)}
}

func (l *logger) WithParams(params Params) Logger {
	return &logger{l.WithFields(logrus.Fields(params))}
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-ID"
var CorrelationIDHeader = "X-Correlation-ID"

// WithRequest returns a context which knows the request ID and correlation ID in the given request.
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	id := getRequestID(req)
	if id == "" {
		id = uuid.New().String()
		req.Header.Set(RequestIDHeader, id)
	}
	ctx = context.WithValue(ctx, requestIDKey, id)
	if id := getCorrelationID(req); id != "" {
		ctx = context.WithValue(ctx, correlationIDKey, id)
	}
	return ctx
}

// GetRequestID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// getCorrelationID extracts the correlation ID from the HTTP request
func getCorrelationID(req *http.Request) string {
	return req.Header.Get(CorrelationIDHeader)
}

// getRequestID extracts the correlation ID from the HTTP request
func getRequestID(req *http.Request) string {
	return req.Header.Get(RequestIDHeader)
}
