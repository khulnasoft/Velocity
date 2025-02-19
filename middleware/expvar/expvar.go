package expvar

import (
	"strings"

	"github.com/valyala/fasthttp/expvarhandler"
	"go.khulnasoft.com/velocity/v3"
)

// New creates a new middleware handler
func New(config ...Config) velocity.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c velocity.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		path := c.Path()
		// We are only interested in /debug/vars routes
		if len(path) < 11 || !strings.HasPrefix(path, "/debug/vars") {
			return c.Next()
		}
		if path == "/debug/vars" {
			expvarhandler.ExpvarHandler(c.RequestCtx())
			return nil
		}

		return c.Redirect().To("/debug/vars")
	}
}
