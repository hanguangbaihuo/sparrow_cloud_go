package opentracing

import (
	"log"
	"strings"

	"github.com/kataras/iris/v12/context"

	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

// ZipkinSpan is span
// if the span-context that extract from incoming request header is not nil, ZipkinSpan will be child span of incoming
// else will be root span
var ZipkinSpan zipkin.Span

// OpentracingInf saving opentracing information that incoming request header
// usually contain "x-request-id", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled"...
// for having installed envoy kubernetes pods, only transmit the header to next pod is enough
var OpentracingInf map[string][]string

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

// ------------Inject zipkin b3 header to reqeust header----------
// if opentracing.GlobalTracer != nil && opentracing.ZipkinSpan != nil {
// 	var operationName string
// 	operationName, ok = kwarg["operationname"]
// 	if !ok {
// 		operationName = destURL
// 	}

// 	appSpan := opentracing.GlobalTracer.StartSpan(operationName,
// 		zipkin.Parent(opentracing.ZipkinSpan.Context()),
// 	)
// 	defer appSpan.Finish()
// 	zipkin.TagHTTPMethod.Set(appSpan, req.Method)
// 	zipkin.TagHTTPPath.Set(appSpan, req.URL.Path)
// 	_ = b3.InjectHTTP(req, b3.WithSingleAndMultiHeader())(appSpan.Context())
// 	// log.Printf("send %s header is %#v\n", destURL, req.Header)
// }

func output(sc model.SpanContext) {
	log.Printf("TraceID: %s\n", sc.TraceID.String())
	if sc.ParentID != nil {
		log.Printf("ParentID: %s\n", (*sc.ParentID).String())
	}
	log.Printf("ID: %s\n", sc.ID.String())
}

// Serve is for saving incoming request header of b3
func Serve(ctx context.Context) {
	OpentracingInf = make(map[string][]string)
	headers := ctx.Request().Header
	for key, value := range headers {
		if strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-") {
			OpentracingInf[key] = value
		}
	}
	ctx.Next()
}
