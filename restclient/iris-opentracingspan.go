package restclient

import (
	"github.com/kataras/iris/v12/context"
	opentracing "github.com/opentracing/opentracing-go"
)

// todo:
// import this variable from opentracing middleware
var activeSpanKey = "myOpenTracingSpan"

// ContextWithSpan returns a new `context.Context` that holds a reference to
// the span. If span is nil, a new context without an active span is returned.
func ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	// if span != nil {
	// 	if tracerWithHook, ok := span.Tracer().(TracerContextWithSpanExtension); ok {
	// 		ctx = tracerWithHook.ContextWithSpanHook(ctx, span)
	// 	}
	// }
	// return context.WithValue(ctx, activeSpanKey, span)
	ctx.Values().Set(activeSpanKey, span)
	return ctx
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
//
// NOTE: context.Context != SpanContext: the former is Go's intra-process
// context propagation mechanism, and the latter houses OpenTracing's per-Span
// identity and baggage information.
func SpanFromContext(ctx context.Context) opentracing.Span {
	val := ctx.Values().Get(activeSpanKey)
	if sp, ok := val.(opentracing.Span); ok {
		return sp
	}
	return nil
}

// StartSpanFromContext starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a ChildOfRef. If no such parent could be
// found, StartSpanFromContext creates a root (parentless) Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return StartSpanFromContextWithTracer(ctx, opentracing.GlobalTracer(), operationName, opts...)
}

// StartSpanFromContextWithTracer starts and returns a span with `operationName`
// using  a span found within the context as a ChildOfRef. If that doesn't exist
// it creates a root span. It also returns a context.Context object built
// around the returned span.
//
// It's behavior is identical to StartSpanFromContext except that it takes an explicit
// tracer as opposed to using the global tracer.
func StartSpanFromContextWithTracer(ctx context.Context, tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(operationName, opts...)
	// ContextWithSpan(ctx, span) function would overwrite opentracing parent span in iris ctx
	// This causes the child span to point to the wrong parent span, so don't change activeSpanKey in ctx
	// It means every service only have two layer: input parent span and output child span
	// return span, ContextWithSpan(ctx, span)
	return span, ctx
}
