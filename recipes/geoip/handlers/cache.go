package handlers

import (
	"time"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cache"
)

// Cache will return a caching middleware
func Cache(exp time.Duration) fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   exp,
		CacheControl: true,
	})
}
