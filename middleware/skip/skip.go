package skip

import (
	"github.com/khulnasoft/velocity"
)

// New creates a middleware handler which skips the wrapped handler
// if the exclude predicate returns true.
func New(handler velocity.Handler, exclude func(c velocity.Ctx) bool) velocity.Handler {
	if exclude == nil {
		return handler
	}

	return func(c velocity.Ctx) error {
		if exclude(c) {
			return c.Next()
		}

		return handler(c)
	}
}
