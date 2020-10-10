package restclient

import (
	spaopentracing "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// This file is from github.com/opentracing/opentracing-go/gocontext.go

// ContextWithSpan set the global parent span in spaopentracing middleware ActiveSpan.
func ContextWithSpan(span opentracing.Span) {
	spaopentracing.ActiveSpan = span
}

// SpanFromContext returns the `Span` previously setted by opentracing middleware, or
// `nil` if no such `Span` could be found.
//
func SpanFromContext() opentracing.Span {
	val := spaopentracing.ActiveSpan
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
func StartSpanFromContext(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return StartSpanFromContextWithTracer(opentracing.GlobalTracer(), operationName, opts...)
}

// StartSpanFromContextWithTracer starts and returns a span with `operationName`
// using  a span found within the context as a ChildOfRef. If that doesn't exist
// it creates a root span. It also returns a context.Context object built
// around the returned span.
//
// It's behavior is identical to StartSpanFromContext except that it takes an explicit
// tracer as opposed to using the global tracer.
func StartSpanFromContextWithTracer(tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if parentSpan := SpanFromContext(); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(operationName, opts...)
	// ContextWithSpan(ctx, span) function would overwrite opentracing parent span in iris ctx
	// This causes the child span to point to the wrong parent span, so don't change activeSpanKey in ctx
	// It means every service only have two layer: input parent span and output child span
	// return span, ContextWithSpan(ctx, span)
	return span
}
