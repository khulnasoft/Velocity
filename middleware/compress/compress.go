package compress

import (
	"github.com/khulnasoft/velocity"
	"github.com/valyala/fasthttp"
)

// New creates a new middleware handler
func New(config ...Config) velocity.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Setup request handlers
	var (
		fctx       = func(_ *fasthttp.RequestCtx) {}
		compressor fasthttp.RequestHandler
	)

	// Setup compression algorithm
	switch cfg.Level {
	case LevelDefault:
		// LevelDefault
		compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
			fasthttp.CompressBrotliDefaultCompression,
			fasthttp.CompressDefaultCompression,
		)
	case LevelBestSpeed:
		// LevelBestSpeed
		compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
			fasthttp.CompressBrotliBestSpeed,
			fasthttp.CompressBestSpeed,
		)
	case LevelBestCompression:
		// LevelBestCompression
		compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
			fasthttp.CompressBrotliBestCompression,
			fasthttp.CompressBestCompression,
		)
	default:
		// LevelDisabled
		return func(c velocity.Ctx) error {
			return c.Next()
		}
	}

	// Return new handler
	return func(c velocity.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Continue stack
		if err := c.Next(); err != nil {
			return err
		}

		// Compress response
		compressor(c.RequestCtx())

		// Return from handler
		return nil
	}
}
