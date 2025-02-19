package keyauth

import (
	"errors"

	"go.khulnasoft.com/velocity/v3"
)

type KeyLookupFunc func(c velocity.Ctx) (string, error)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip middleware.
	// Optional. Default: nil
	Next func(velocity.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid key.
	// Optional. Default: nil
	SuccessHandler velocity.Handler

	// ErrorHandler defines a function which is executed for an invalid key.
	// It may be used to define a custom error.
	// Optional. Default: 401 Invalid or expired key
	ErrorHandler velocity.ErrorHandler

	CustomKeyLookup KeyLookupFunc

	// Validator is a function to validate key.
	Validator func(velocity.Ctx, string) (bool, error)

	// KeyLookup is a string in the form of "<source>:<name>" that is used
	// to extract key from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "form:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	KeyLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default value "Bearer".
	AuthScheme string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	SuccessHandler: func(c velocity.Ctx) error {
		return c.Next()
	},
	ErrorHandler: func(c velocity.Ctx, err error) error {
		if errors.Is(err, ErrMissingOrMalformedAPIKey) {
			return c.Status(velocity.StatusUnauthorized).SendString(err.Error())
		}
		return c.Status(velocity.StatusUnauthorized).SendString("Invalid or expired API Key")
	},
	KeyLookup:       "header:" + velocity.HeaderAuthorization,
	CustomKeyLookup: nil,
	AuthScheme:      "Bearer",
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
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = ConfigDefault.SuccessHandler
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = ConfigDefault.ErrorHandler
	}
	if cfg.KeyLookup == "" {
		cfg.KeyLookup = ConfigDefault.KeyLookup
		// set AuthScheme as "Bearer" only if KeyLookup is set to default.
		if cfg.AuthScheme == "" {
			cfg.AuthScheme = ConfigDefault.AuthScheme
		}
	}
	if cfg.Validator == nil {
		panic("velocity: keyauth middleware requires a validator function")
	}

	return cfg
}
