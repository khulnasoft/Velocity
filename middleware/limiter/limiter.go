package limiter

import (
	"github.com/khulnasoft/velocity"
)

const (
	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

type Handler interface {
	New(config Config) velocity.Handler
}

// New creates a new middleware handler
func New(config ...Config) velocity.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return the specified middleware handler.
	return cfg.LimiterMiddleware.New(cfg)
}
