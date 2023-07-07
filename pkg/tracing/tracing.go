package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/star-table/usercenter/core/conf"
)

func EnableTracing() bool {
	return conf.Cfg.Jaeger != nil
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return opentracing.GlobalTracer().StartSpan(operationName, opts...)
}

func Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return opentracing.GlobalTracer().Inject(sm, format, carrier)
}

func Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return opentracing.GlobalTracer().Extract(format, carrier)
}
