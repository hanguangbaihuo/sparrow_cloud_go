package opentracing

import (
	"strings"
	"sync"

	"github.com/kataras/iris/v12/context"
)

// OpentracingHeader is saving opentracing b3 header data
// field Headers is map[string][]string type
type OpentracingHeader struct {
	sync.RWMutex
	Headers interface{}
}

// OpentracingInf global variable to transmit opentracing b3 header data
var OpentracingInf = new(OpentracingHeader)

// SetHeaders set OpentracingHeader Headers function
func (oh *OpentracingHeader) SetHeaders(data interface{}) {
	oh.Lock()
	oh.Headers = data
	oh.Unlock()
}

// GetHeaders get OpentracingHeader Headers function
func (oh *OpentracingHeader) GetHeaders() interface{} {
	oh.RLock()
	data := oh.Headers
	oh.RUnlock()
	return data
}

// Serve is for saving incoming request header of b3
func Serve(ctx context.Context) {
	data := make(map[string][]string)
	headers := ctx.Request().Header
	for key, value := range headers {
		if strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-") {
			data[key] = value
		}
	}
	OpentracingInf.SetHeaders(data)
	ctx.Next()
}

// Inject b3 header to req header
// b3headers, ok := opentracing.OpentracingInf.GetHeaders().(map[string][]string)
// if ok {
// 	for key, value := range b3headers {
// 		if len(value) > 0 {
// 			req.Header.Set(key, value[0])
// 		}
// 	}
// }
