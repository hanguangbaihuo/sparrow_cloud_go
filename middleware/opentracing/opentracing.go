package opentracing

import (
	"log"

	"github.com/kataras/iris/v12/context"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var ActiveSpanKey = "myOpenTracingSpan"

// Serve is opentracing middleware function
// parameter operationName is a operaion name, usually named it service name
func Serve(operationName string) func(context.Context) {
	return func(ctx context.Context) {
		var serverSpan opentracing.Span
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(ctx.Request().Header))
		if err != nil {
			log.Printf("extract global tracer error: %v\n", err)
		}

		serverSpan = opentracing.StartSpan(
			operationName,
			ext.RPCServerOption(wireContext))
		defer serverSpan.Finish()

		ctx.Values().Set(ActiveSpanKey, serverSpan)

		ctx.Next()
	}
}
