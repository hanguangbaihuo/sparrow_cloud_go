package opentracing

import (
	"strings"

	"github.com/kataras/iris/v12/context"
)

// Entry is key value store struct
type Entry struct {
	Key   string
	Value []string
}

// HeaderStore is global variable to saving opentracing b3 headers
var HeaderStore []Entry

// ServeBk is for saving incoming request header of b3
// it is a backup fucntion for Serve in this middleware
func ServeBk(ctx context.Context) {
	// reset HeaderStore for each incoming request
	HeaderStore = HeaderStore[0:0]
	headers := ctx.Request().Header
	for key, value := range headers {
		if strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-") {
			HeaderStore = append(HeaderStore, Entry{key, value})
		}
	}
	ctx.Next()
}

// Inject b3 header to req header
// for _, entry := range opentracing.HeaderStore {
// 	key, value := entry.Key, entry.Value
// 	if len(value) > 0 {
// 		req.Header.Set(key, value[0])
// 	}
// }
