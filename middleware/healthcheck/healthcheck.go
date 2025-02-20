package healthcheck

import (
	"github.com/khulnasoft/velocity"
)

// HealthChecker defines a function to check liveness or readiness of the application
type HealthChecker func(velocity.Ctx) bool

func NewHealthChecker(config ...Config) velocity.Handler {
	cfg := defaultConfigV3(config...)

	return func(c velocity.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if c.Method() != velocity.MethodGet {
			return c.Next()
		}

		if cfg.Probe(c) {
			return c.SendStatus(velocity.StatusOK)
		}

		return c.SendStatus(velocity.StatusServiceUnavailable)
	}
}
