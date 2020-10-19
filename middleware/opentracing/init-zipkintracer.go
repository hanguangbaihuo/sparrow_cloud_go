package opentracing

import (
	"log"

	zipkin "github.com/openzipkin/zipkin-go"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

// GlobalTracer is zipkin Tracer
var GlobalTracer *zipkin.Tracer
var zipkinConfig Config

// Config is zipkin Global Tracer configuration
type Config struct {
	// Debug configure zipkin output
	// default false
	Debug bool
}

// InitZipkinTracer is for initilize zipkin global Tracer function
func InitZipkinTracer(debug bool) *zipkin.Tracer {
	var err error
	// set up zipkin config
	zipkinConfig = Config{Debug: debug}
	if debug {
		// set up a span reporter
		// this reporter can be replaced with amqp, http, kafka, recorder...
		reporter := logreporter.NewReporter(nil)
		defer reporter.Close()
		// initialize our tracer
		GlobalTracer, err = zipkin.NewTracer(reporter)
	} else {
		// initialize our tracer
		GlobalTracer, err = zipkin.NewTracer(nil)
	}
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	return GlobalTracer
}
