package opentracing

import (
	"log"

	"github.com/kataras/iris/v12/context"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

// ActiveSpan is global opentraing parent span, will be initialized in opentracing Serve function
// restclient use it to generate child span
var ActiveSpan opentracing.Span

// ZipkinSpan is span
// if the span-context that extract from incoming request header is not nil, ZipkinSpan will be child span of incoming
// else will be root span
var ZipkinSpan zipkin.Span

// ZipkinSpanContext is span context from income request header
// var ZipkinSpanContext model.SpanContext

// Serve is opentracing middleware function
// parameter operationName is a operaion name, usually named it service name
func Serve(operationName string) func(context.Context) {
	return func(ctx context.Context) {
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(ctx.Request().Header))
		if err != nil {
			log.Printf("global tracer extract span context not found: %v\n", err)
		}

		ActiveSpan = opentracing.StartSpan(
			operationName,
			ext.RPCServerOption(wireContext))
		defer ActiveSpan.Finish()

		// ctx.Values().Set(ActiveSpanKey, ActiveSpan)

		ctx.Next()
	}
}

func ZipkinServe(operationName string) func(context.Context) {
	return func(ctx context.Context) {
		if GlobalTracer == nil {
			log.Printf("[WARNNING] Before using Zipkin opentracing middleware, PLZ init InitZipkinTracer, otherwise it do NOT work!!!\n")
			ctx.Next()
			return
		}
		sc := GlobalTracer.Extract(b3.ExtractHTTP(ctx.Request()))
		if sc.Err != nil {
			log.Printf("global tracer extract span context not found: %v\n", sc.Err)
		}

		ZipkinSpan = GlobalTracer.StartSpan(operationName,
			// zipkin.Kind(model.Server),
			zipkin.Parent(sc),
		)
		// defer ZipkinSpan.Finish()

		// todo: add tag to ZipkinSpan
		if zipkinConfig.Debug {
			log.Printf("------------------ origin span context ------------------")
			output(sc)
			log.Printf("------------------ generate span context ------------------")
			output(ZipkinSpan.Context())
		}
		ctx.Next()
	}
}

func output(sc model.SpanContext) {
	log.Printf("TraceID: %s\n", sc.TraceID.String())
	if sc.ParentID != nil {
		log.Printf("ParentID: %s\n", (*sc.ParentID).String())
	}
	log.Printf("ID: %s\n", sc.ID.String())
}
