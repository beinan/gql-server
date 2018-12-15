package logging

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitOpenTracing(serviceName string) io.Closer {
	tracer, closer, err := initJaegerTracer(serviceName)
	if err != nil {

	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}

func initJaegerTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: "jaeger:5778",
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}
	return cfg.NewTracer()
}

func StartSpanFromContext(
	ctx context.Context, name string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, name)
}
