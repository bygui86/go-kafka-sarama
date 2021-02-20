package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	zipkinopentracing "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttpreporter "github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
)

const (
	zipkinUrlFormat       = "http://%s:%d/api/v2/spans"
	zipkinDefaultHostPort = ":0"
)

func InitSampleZipkin(serviceName, zipkinHost string, zipkinPort int) (reporter.Reporter, error) {
	// set up a span reporter
	reporter := zipkinhttpreporter.NewReporter(fmt.Sprintf(zipkinUrlFormat, zipkinHost, zipkinPort))

	// create our local service endpoint
	endpoint, endpointErr := zipkin.NewEndpoint(serviceName, zipkinDefaultHostPort)
	if endpointErr != nil {
		return nil, endpointErr
	}

	// initialize our tracer
	nativeTracer, tracerErr := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
	)
	if tracerErr != nil {
		return nil, tracerErr
	}

	// use zipkin-go-opentracing to wrap Zipkin tracer into OpenTracing tracer
	// and set it as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(zipkinopentracing.Wrap(nativeTracer))

	logging.SugaredLog.Debugf("Zipkin global tracer registered: %t", opentracing.IsGlobalTracerRegistered())
	return reporter, nil
}

func InitTestingZipkin(serviceName, zipkinHost string, zipkinPort int) (reporter.Reporter, error) {
	// set up a span reporter
	reporter := zipkinhttpreporter.NewReporter(fmt.Sprintf(zipkinUrlFormat, zipkinHost, zipkinPort))

	// create our local service endpoint
	endpoint, endpointErr := zipkin.NewEndpoint(serviceName, zipkinDefaultHostPort)
	if endpointErr != nil {
		return nil, endpointErr
	}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	// initialize our tracer
	nativeTracer, tracerErr := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(endpoint),
	)
	if tracerErr != nil {
		return nil, tracerErr
	}

	// use zipkin-go-opentracing to wrap Zipkin tracer into OpenTracing tracer
	// and set it as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(zipkinopentracing.Wrap(nativeTracer))

	logging.SugaredLog.Debugf("Zipkin global tracer registered: %t", opentracing.IsGlobalTracerRegistered())
	return reporter, nil
}

func InitProductionZipkin(serviceName, zipkinHost string, zipkinPort int) (reporter.Reporter, error) {
	// set up a span reporter
	reporter := zipkinhttpreporter.NewReporter(fmt.Sprintf(zipkinUrlFormat, zipkinHost, zipkinPort))

	// create our local service endpoint
	endpoint, endpointErr := zipkin.NewEndpoint(serviceName, zipkinDefaultHostPort)
	if endpointErr != nil {
		return nil, endpointErr
	}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(0.2)
	if err != nil {
		return nil, err
	}

	// initialize our tracer
	nativeTracer, tracerErr := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(endpoint),
	)
	if tracerErr != nil {
		return nil, tracerErr
	}

	// use zipkin-go-opentracing to wrap Zipkin tracer into OpenTracing tracer
	// and set it as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(zipkinopentracing.Wrap(nativeTracer))

	logging.SugaredLog.Debugf("Zipkin global tracer registered: %t", opentracing.IsGlobalTracerRegistered())
	return reporter, nil
}
