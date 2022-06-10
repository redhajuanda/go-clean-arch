package otel

import (
	"context"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string, service string, version string, env string, sampled bool) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	var sampler = tracesdk.NeverSample()
	if sampled {
		sampler = tracesdk.AlwaysSample()
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", env),
		)),
	)
	return tp, nil
}

func SetTraceProvider(url string, service string, version string, env string, sampled bool) error {
	tp, err := tracerProvider(url, service, version, env, sampled)
	if err != nil {
		return err
	}
	otel.SetTracerProvider(tp)
	return nil
}

func Start(ctx context.Context) (context.Context, trace.Span) {

	c, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(c).Name()
	fs := strings.SplitN(f, ".", 2)
	replacer := strings.NewReplacer("(", "", ")", "", "*", "")
	operation := replacer.Replace(fs[1])
	return otel.Tracer(fs[0]).Start(ctx, operation)
}
