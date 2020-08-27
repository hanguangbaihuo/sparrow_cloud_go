package accesscontrol

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var (
	// ErrResrouceMissing 未提供资源
	ErrResrouceMissing = errors.New("required resource not found")
)

// // A function called whenever an error is encountered
type errorHandler func(context.Context, error)

// Config is
type Config struct {
	// 资源
	ResourceValue string
	// The function that will be called when there's an error validating the token
	// Default value:
	ErrorHandler errorHandler
}

// Middleware the middleware for JSON Web tokens authentication method
type Middleware struct {
	Config Config
}

// OnError is the default error handler.
// Use it to change the behavior for each error.
// See `Config.ErrorHandler`.
func OnError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	ctx.StatusCode(iris.StatusForbidden)
	ctx.WriteString(err.Error())
}

// New constructs a new Secure instance with supplied options.
func New(resource string) *Middleware {

	var c Config
	c = Config{}
	c.ResourceValue = resource

	if c.ErrorHandler == nil {
		c.ErrorHandler = OnError
	}

	return &Middleware{Config: c}
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx context.Context) {
	if err := m.CheckAccessControl(ctx); err != nil {
		m.Config.ErrorHandler(ctx, err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// CheckAccessControl the main functionality, checks for token
func (m *Middleware) CheckAccessControl(ctx context.Context) error {
	resource := m.Config.ResourceValue
	if len(resource) == 0 {
		return ErrResrouceMissing
	}
	return nil
}
