package earlydata

import (
	"github.com/khulnasoft/velocity"
)

const (
	DefaultHeaderName      = "Early-Data"
	DefaultHeaderTrueValue = "1"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c velocity.Ctx) bool

	// IsEarlyData returns whether the request is an early-data request.
	//
	// Optional. Default: a function which checks if the "Early-Data" request header equals "1".
	IsEarlyData func(c velocity.Ctx) bool

	// AllowEarlyData returns whether the early-data request should be allowed or rejected.
	//
	// Optional. Default: a function which rejects the request on unsafe and allows the request on safe HTTP request methods.
	AllowEarlyData func(c velocity.Ctx) bool

	// Error is returned in case an early-data request is rejected.
	//
	// Optional. Default: velocity.ErrTooEarly.
	Error error
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	IsEarlyData: func(c velocity.Ctx) bool {
		return c.Get(DefaultHeaderName) == DefaultHeaderTrueValue
	},

	AllowEarlyData: func(c velocity.Ctx) bool {
		return velocity.IsMethodSafe(c.Method())
	},

	Error: velocity.ErrTooEarly,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values

	if cfg.IsEarlyData == nil {
		cfg.IsEarlyData = ConfigDefault.IsEarlyData
	}

	if cfg.AllowEarlyData == nil {
		cfg.AllowEarlyData = ConfigDefault.AllowEarlyData
	}

	if cfg.Error == nil {
		cfg.Error = ConfigDefault.Error
	}

	return cfg
}
