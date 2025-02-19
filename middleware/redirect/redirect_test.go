package redirect

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity/v3"
)

func Test_Redirect(t *testing.T) {
	app := *velocity.New()

	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/default/*": "velocity.wiki",
		},
		StatusCode: velocity.StatusTemporaryRedirect,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/redirect/*": "$1",
		},
		StatusCode: velocity.StatusSeeOther,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/pattern/*": "golang.org",
		},
		StatusCode: velocity.StatusFound,
	}))

	app.Use(New(Config{
		Rules: map[string]string{
			"/": "/swagger",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/params": "/with_params",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))

	app.Get("/api/*", func(c velocity.Ctx) error {
		return c.SendString("API")
	})

	app.Get("/new", func(c velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	tests := []struct {
		name       string
		url        string
		redirectTo string
		statusCode int
	}{
		{
			name:       "should be returns status StatusFound without a wildcard",
			url:        "/default",
			redirectTo: "google.com",
			statusCode: velocity.StatusMovedPermanently,
		},
		{
			name:       "should be returns status StatusTemporaryRedirect  using wildcard",
			url:        "/default/xyz",
			redirectTo: "velocity.wiki",
			statusCode: velocity.StatusTemporaryRedirect,
		},
		{
			name:       "should be returns status StatusSeeOther without set redirectTo to use the default",
			url:        "/redirect/github.com/khulnasoft/redirect",
			redirectTo: "github.com/khulnasoft/redirect",
			statusCode: velocity.StatusSeeOther,
		},
		{
			name:       "should return the status code default",
			url:        "/pattern/xyz",
			redirectTo: "golang.org",
			statusCode: velocity.StatusFound,
		},
		{
			name:       "access URL without rule",
			url:        "/new",
			statusCode: velocity.StatusOK,
		},
		{
			name:       "redirect to swagger route",
			url:        "/",
			redirectTo: "/swagger",
			statusCode: velocity.StatusMovedPermanently,
		},
		{
			name:       "no redirect to swagger route",
			url:        "/api/",
			statusCode: velocity.StatusOK,
		},
		{
			name:       "no redirect to swagger route #2",
			url:        "/api/test",
			statusCode: velocity.StatusOK,
		},
		{
			name:       "redirect with query params",
			url:        "/params?query=abc",
			redirectTo: "/with_params?query=abc",
			statusCode: velocity.StatusMovedPermanently,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), velocity.MethodGet, tt.url, nil)
			require.NoError(t, err)
			req.Header.Set("Location", "github.com/khulnasoft/redirect")
			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tt.statusCode, resp.StatusCode)
			require.Equal(t, tt.redirectTo, resp.Header.Get("Location"))
		})
	}
}

func Test_Next(t *testing.T) {
	// Case 1 : Next function always returns true
	app := *velocity.New()
	app.Use(New(Config{
		Next: func(velocity.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))

	app.Use(func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, velocity.StatusOK, resp.StatusCode)

	// Case 2 : Next function always returns false
	app = *velocity.New()
	app.Use(New(Config{
		Next: func(velocity.Ctx) bool {
			return false
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)

	require.Equal(t, velocity.StatusMovedPermanently, resp.StatusCode)
	require.Equal(t, "google.com", resp.Header.Get("Location"))
}

func Test_NoRules(t *testing.T) {
	// Case 1: No rules with default route defined
	app := *velocity.New()

	app.Use(New(Config{
		StatusCode: velocity.StatusMovedPermanently,
	}))

	app.Use(func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)

	// Case 2: No rules and no default route defined
	app = *velocity.New()

	app.Use(New(Config{
		StatusCode: velocity.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
}

func Test_DefaultConfig(t *testing.T) {
	// Case 1: Default config and no default route
	app := *velocity.New()

	app.Use(New())

	req, err := http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)

	// Case 2: Default config and default route
	app = *velocity.New()

	app.Use(New())
	app.Use(func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)

	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
}

func Test_RegexRules(t *testing.T) {
	// Case 1: Rules regex is empty
	app := *velocity.New()
	app.Use(New(Config{
		Rules:      map[string]string{},
		StatusCode: velocity.StatusMovedPermanently,
	}))

	app.Use(func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)

	// Case 2: Rules regex map contains valid regex and well-formed replacement URLs
	app = *velocity.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: velocity.StatusMovedPermanently,
	}))

	app.Use(func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), velocity.MethodGet, "/default", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)

	require.NoError(t, err)
	require.Equal(t, velocity.StatusMovedPermanently, resp.StatusCode)
	require.Equal(t, "google.com", resp.Header.Get("Location"))

	// Case 3: Test invalid regex throws panic
	app = *velocity.New()
	require.Panics(t, func() {
		app.Use(New(Config{
			Rules: map[string]string{
				"(": "google.com",
			},
			StatusCode: velocity.StatusMovedPermanently,
		}))
	})
}
