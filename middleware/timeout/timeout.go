package timeout

import (
	"context"
	"errors"
	"time"

	"go.khulnasoft.com/velocity/v3"
)

// New enforces a timeout for each incoming request. If the timeout expires or
// any of the specified errors occur, velocity.ErrRequestTimeout is returned.
func New(h velocity.Handler, timeout time.Duration, tErrs ...error) velocity.Handler {
	return func(ctx velocity.Ctx) error {
		// If timeout <= 0, skip context.WithTimeout and run the handler as-is.
		if timeout <= 0 {
			return runHandler(ctx, h, tErrs)
		}

		// Create a context with the specified timeout; any operation exceeding
		// this deadline will be canceled automatically.
		timeoutContext, cancel := context.WithTimeout(ctx.Context(), timeout)
		defer cancel()

		// Replace the default Velocity context with our timeout-bound context.
		ctx.SetContext(timeoutContext)

		// Run the handler and check for relevant errors.
		err := runHandler(ctx, h, tErrs)

		// If the context actually timed out, return a timeout error.
		if errors.Is(timeoutContext.Err(), context.DeadlineExceeded) {
			return velocity.ErrRequestTimeout
		}
		return err
	}
}

// runHandler executes the handler and returns velocity.ErrRequestTimeout if it
// sees a deadline exceeded error or one of the custom "timeout-like" errors.
func runHandler(c velocity.Ctx, h velocity.Handler, tErrs []error) error {
	// Execute the wrapped handler synchronously.
	err := h(c)
	// If the context has timed out, return a request timeout error.
	if err != nil && (errors.Is(err, context.DeadlineExceeded) || isCustomError(err, tErrs)) {
		return velocity.ErrRequestTimeout
	}
	return err
}

// isCustomError checks whether err matches any error in errList using errors.Is.
func isCustomError(err error, errList []error) bool {
	for _, e := range errList {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}
