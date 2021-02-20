package tracing

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opentracing/opentracing-go"
)

/*
	KafkaHeadersCarrier satisfies both TextMapWriter and TextMapReader.
	See Inject(..) method as example usage on producer side.
	See Extract(..) method as example usage on consumer side.
*/
type KafkaHeadersCarrier []kafka.Header

// Set implements Set() of opentracing.TextMapWriter interface.
func (c *KafkaHeadersCarrier) Set(key, val string) {
	*c = append(*c, kafka.Header{
		Key:   key,
		Value: []byte(val),
	})
}

// ForeachKey implements ForeachKey(..) of opentracing.TextMapReader interface.
func (c *KafkaHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for _, val := range *c {
		if val.Key == "" && string(val.Value) == "" {
			continue
		}
		if err := handler(val.Key, string(val.Value)); err != nil {
			return err
		}
	}
	return nil
}

// Inject injects the span context into Kafka headers.
func Inject(span opentracing.Span, carrier *KafkaHeadersCarrier) error {
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	// ALTERNATIVE with specific Tracer
	// err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
	return err
}

// Extract extracts the span context out of Kafka headers.
func Extract(carrier *KafkaHeadersCarrier) (opentracing.SpanContext, error) {
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)
}
