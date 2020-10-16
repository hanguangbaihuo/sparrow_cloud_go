package opentracing

import (
	"fmt"
	"io"
	"log"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaezipkin "github.com/uber/jaeger-client-go/zipkin"

	zipkin "github.com/openzipkin/zipkin-go"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

var GlobalTracer *zipkin.Tracer

//InitGlobalTracer is for setting Global Tracer.
// you must do the function in main() function like follow:
//   closer := InitGlobalTracer("your_service_name")
//   defer closer.Close()
func InitGlobalTracer(serviceName string) io.Closer {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
	zipkinPropagator := jaezipkin.NewZipkinB3HTTPHeaderPropagator()

	// Initialize tracer with a logger and a metrics factory
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	closer, err := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)

	if err != nil {
		panic(fmt.Sprintf("ERROR: init Jaeger Opentracing: %v\n", err))
	}

	return closer
}

func InitZipkinTracer(serviceName string) *zipkin.Tracer {
	// set up a span reporter
	reporter := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	defer reporter.Close()

	// initialize our tracer
	GlobalTracer, err := zipkin.NewTracer(reporter)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	return GlobalTracer
}
