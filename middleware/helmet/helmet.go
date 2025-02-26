package helmet

import (
	"fmt"

	"github.com/khulnasoft/velocity"
)

// New creates a new middleware handler
func New(config ...Config) velocity.Handler {
	// Init config
	cfg := configDefault(config...)

	// Return middleware handler
	return func(c velocity.Ctx) error {
		// Next request to skip middleware
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Set headers
		if cfg.XSSProtection != "" {
			c.Set(velocity.HeaderXXSSProtection, cfg.XSSProtection)
		}

		if cfg.ContentTypeNosniff != "" {
			c.Set(velocity.HeaderXContentTypeOptions, cfg.ContentTypeNosniff)
		}

		if cfg.XFrameOptions != "" {
			c.Set(velocity.HeaderXFrameOptions, cfg.XFrameOptions)
		}

		if cfg.CrossOriginEmbedderPolicy != "" {
			c.Set("Cross-Origin-Embedder-Policy", cfg.CrossOriginEmbedderPolicy)
		}

		if cfg.CrossOriginOpenerPolicy != "" {
			c.Set("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
		}

		if cfg.CrossOriginResourcePolicy != "" {
			c.Set("Cross-Origin-Resource-Policy", cfg.CrossOriginResourcePolicy)
		}

		if cfg.OriginAgentCluster != "" {
			c.Set("Origin-Agent-Cluster", cfg.OriginAgentCluster)
		}

		if cfg.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", cfg.ReferrerPolicy)
		}

		if cfg.XDNSPrefetchControl != "" {
			c.Set("X-DNS-Prefetch-Control", cfg.XDNSPrefetchControl)
		}

		if cfg.XDownloadOptions != "" {
			c.Set("X-Download-Options", cfg.XDownloadOptions)
		}

		if cfg.XPermittedCrossDomain != "" {
			c.Set("X-Permitted-Cross-Domain-Policies", cfg.XPermittedCrossDomain)
		}

		// Handle HSTS headers
		if c.Protocol() == "https" && cfg.HSTSMaxAge != 0 {
			subdomains := ""
			if !cfg.HSTSExcludeSubdomains {
				subdomains = "; includeSubDomains"
			}
			if cfg.HSTSPreloadEnabled {
				subdomains += "; preload"
			}
			c.Set(velocity.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
		}

		// Handle Content-Security-Policy headers
		if cfg.ContentSecurityPolicy != "" {
			if cfg.CSPReportOnly {
				c.Set(velocity.HeaderContentSecurityPolicyReportOnly, cfg.ContentSecurityPolicy)
			} else {
				c.Set(velocity.HeaderContentSecurityPolicy, cfg.ContentSecurityPolicy)
			}
		}

		// Handle Permissions-Policy headers
		if cfg.PermissionPolicy != "" {
			c.Set(velocity.HeaderPermissionsPolicy, cfg.PermissionPolicy)
		}

		return c.Next()
	}
}
