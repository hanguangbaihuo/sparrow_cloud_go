package opentracing

import (
	"log"

	"github.com/kataras/iris/v12/context"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// ActiveSpan is global opentraing parent span, will be initialized in opentracing Serve function
// restclient use it to generate child span
var ActiveSpan opentracing.Span

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
