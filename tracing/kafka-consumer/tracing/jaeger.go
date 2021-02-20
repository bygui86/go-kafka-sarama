package tracing

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLogZap "github.com/uber/jaeger-client-go/log/zap"
	jaegerProm "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
)

/*
	Environment variables:
		JAEGER_DISABLED to enable/disable Jaeger Tracer (default false)
		JAEGER_AGENT_HOST to set target Jaeger host (default localhost)
		JAEGER_AGENT_PORT to set target Jaeger port (default 6831)
		JAEGER_SAMPLER_TYPE to set Sampler type (default none)(*)
		JAEGER_SAMPLER_PARAM to set Sampler param (default 0) (**)
		JAEGER_REPORTER_MAX_QUEUE_SIZE to set max how many spans the reporter can keep in memory before it starts
			dropping new spans (default 0). The queue is continuously drained by a background go-routine, as
			fast as spans can be sent out of process.
		JAEGER_REPORTER_LOG_SPANS to set Reporter LogSpans (default false)
		JAEGER_REPORTER_FLUSH_INTERVAL to set Reporter flush interval (default 0s) (***)
		JAEGER_SERVICE_NAME	to set service name on the side of Jaeger (default empty string)

	(*) JAEGER_SAMPLER_TYPE available values:
			"const" (SamplerTypeConst) is the type of sampler that always makes the same decision.
			"remote" (SamplerTypeRemote) is the type of sampler that polls Jaeger agent for sampling strategy.
			"probabilistic" (SamplerTypeProbabilistic) is the type of sampler that samples traces with a certain fixed probability.
			"ratelimiting" (SamplerTypeRateLimiting) is the type of sampler that samples only up to a fixed number of traces per second.

	(**) JAEGER_SAMPLER_PARAM available values:
			for "const" sampler, 0 or 1 for always false/true respectively
			for "probabilistic" sampler, a probability between 0 and 1
			for "rateLimiting" sampler, the number of spans per second
			for "remote" sampler, param is the same as for "probabilistic"
				and indicates the initial sampling rate before the actual one
				is received from the mothership.

	(***) JAEGER_REPORTER_FLUSH_INTERVAL valid time units:
			"ns", "us" (or "µs"), "ms", "s", "m", "h".
*/

/*
	Use Zap logger to print out spans.
	Use a Prometheus registerer to expose metrics.
*/
func InitTracer() (io.Closer, error) {
	cfg, cfgErr := jaegerCfg.FromEnv()
	if cfgErr != nil {
		return nil, cfgErr
	}
	logging.SugaredLog.Debugf("Jaeger Configuration: %+v", cfg)
	logging.SugaredLog.Debugf("Jaeger Sampler: %+v", cfg.Sampler)
	logging.SugaredLog.Debugf("Jaeger Reporter: %+v", cfg.Reporter)

	closer, tracerErr := cfg.InitGlobalTracer(
		cfg.ServiceName,
		jaegerCfg.Logger(jaegerLogZap.NewLogger(logging.Log)),
		jaegerCfg.Metrics(jaegerProm.New(jaegerProm.WithRegisterer(prometheus.DefaultRegisterer))),
	)
	if tracerErr != nil {
		return nil, tracerErr
	}

	logging.SugaredLog.Debugf("Jaeger global Tracer registered: %t", opentracing.IsGlobalTracerRegistered())
	return closer, nil
}

/*
	Sample configuration for testing. Use constant sampling to sample every trace and enable LogSpan to log every
	span via configured Logger. Use a Prometheus registerer to expose metrics.
*/
func SampleManualInitTracer(serviceName string) (io.Closer, error) {
	cfg := jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	closer, tracerErr := cfg.InitGlobalTracer(
		serviceName,
		jaegerCfg.Logger(jaegerLogZap.NewLogger(logging.Log)),
		jaegerCfg.Metrics(jaegerProm.New(jaegerProm.WithRegisterer(prometheus.DefaultRegisterer))),
	)
	if tracerErr != nil {
		return nil, tracerErr
	}

	logging.SugaredLog.Debugf("Jaeger global tracer registered: %t", opentracing.IsGlobalTracerRegistered())
	return closer, nil
}
