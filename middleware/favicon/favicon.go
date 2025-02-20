package favicon

import (
	"io"
	"io/fs"
	"os"
	"strconv"

	"go.khulnasoft.com/velocity"
)

// Config defines the config for middleware.
type Config struct {
	// FileSystem is an optional alternate filesystem to search for the favicon in.
	// An example of this could be an embedded or network filesystem
	//
	// Optional. Default: nil
	FileSystem fs.FS `json:"-"`

	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c velocity.Ctx) bool

	// File holds the path to an actual favicon that will be cached
	//
	// Optional. Default: ""
	File string `json:"file"`

	// URL for favicon handler
	//
	// Optional. Default: "/favicon.ico"
	URL string `json:"url"`

	// CacheControl defines how the Cache-Control header in the response should be set
	//
	// Optional. Default: "public, max-age=31536000"
	CacheControl string `json:"cache_control"`

	// Raw data of the favicon file
	//
	// Optional. Default: nil
	Data []byte `json:"-"`
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:         nil,
	File:         "",
	URL:          fPath,
	CacheControl: "public, max-age=31536000",
}

const (
	fPath  = "/favicon.ico"
	hType  = "image/x-icon"
	hAllow = "GET, HEAD, OPTIONS"
	hZero  = "0"
)

// New creates a new middleware handler
func New(config ...Config) velocity.Handler {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Next == nil {
			cfg.Next = ConfigDefault.Next
		}
		if cfg.URL == "" {
			cfg.URL = ConfigDefault.URL
		}
		if cfg.File == "" {
			cfg.File = ConfigDefault.File
		}
		if cfg.CacheControl == "" {
			cfg.CacheControl = ConfigDefault.CacheControl
		}
	}

	// Load iconData if provided
	var (
		err           error
		iconData      []byte
		iconLenHeader string
		iconLen       int
	)
	if cfg.Data != nil {
		// use the provided favicon data
		iconData = cfg.Data
		iconLenHeader = strconv.Itoa(len(cfg.Data))
		iconLen = len(cfg.Data)
	} else if cfg.File != "" {
		// read from configured filesystem if present
		if cfg.FileSystem != nil {
			f, err := cfg.FileSystem.Open(cfg.File)
			if err != nil {
				panic(err)
			}
			if iconData, err = io.ReadAll(f); err != nil {
				panic(err)
			}
		} else if iconData, err = os.ReadFile(cfg.File); err != nil {
			panic(err)
		}

		iconLenHeader = strconv.Itoa(len(iconData))
		iconLen = len(iconData)
	}

	// Return new handler
	return func(c velocity.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Only respond to favicon requests
		if c.Path() != cfg.URL {
			return c.Next()
		}

		// Only allow GET, HEAD and OPTIONS requests
		if c.Method() != velocity.MethodGet && c.Method() != velocity.MethodHead {
			if c.Method() != velocity.MethodOptions {
				c.Status(velocity.StatusMethodNotAllowed)
			} else {
				c.Status(velocity.StatusOK)
			}
			c.Set(velocity.HeaderAllow, hAllow)
			c.Set(velocity.HeaderContentLength, hZero)
			return nil
		}

		// Serve cached favicon
		if iconLen > 0 {
			c.Set(velocity.HeaderContentLength, iconLenHeader)
			c.Set(velocity.HeaderContentType, hType)
			c.Set(velocity.HeaderCacheControl, cfg.CacheControl)
			return c.Status(velocity.StatusOK).Send(iconData)
		}

		return c.SendStatus(velocity.StatusNoContent)
	}
}
