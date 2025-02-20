package earlydata_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/earlydata"
)

const (
	headerName   = "Early-Data"
	headerValOn  = "1"
	headerValOff = "0"
)

func appWithConfig(t *testing.T, c *velocity.Config) *velocity.App {
	t.Helper()
	t.Parallel()

	var app *velocity.App
	if c == nil {
		app = velocity.New()
	} else {
		app = velocity.New(*c)
	}

	app.Use(earlydata.New())

	// Middleware to test IsEarly func
	const localsKeyTestValid = "earlydata_testvalid"
	app.Use(func(c velocity.Ctx) error {
		isEarly := earlydata.IsEarly(c)

		switch h := c.Get(headerName); h {
		case "", headerValOff:
			if isEarly {
				return errors.New("is early-data even though it's not")
			}

		case headerValOn:
			switch {
			case velocity.IsMethodSafe(c.Method()):
				if !isEarly {
					return errors.New("should be early-data on safe HTTP methods")
				}
			default:
				if isEarly {
					return errors.New("early-data unsuported on unsafe HTTP methods")
				}
			}

		default:
			return fmt.Errorf("header has unsupported value: %s", h)
		}

		_ = c.Locals(localsKeyTestValid, true)

		return c.Next()
	})

	app.Add([]string{
		velocity.MethodGet,
		velocity.MethodPost,
	}, "/", func(c velocity.Ctx) error {
		valid, ok := c.Locals(localsKeyTestValid).(bool)
		if !ok {
			panic(errors.New("failed to type-assert to bool"))
		}
		if !valid {
			return errors.New("handler called even though validation failed")
		}

		return nil
	})

	return app
}

// go test -run Test_EarlyData
func Test_EarlyData(t *testing.T) {
	t.Parallel()

	trustedRun := func(t *testing.T, app *velocity.App) {
		t.Helper()

		{
			req := httptest.NewRequest(velocity.MethodGet, "/", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)

			req.Header.Set(headerName, headerValOff)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)

			req.Header.Set(headerName, headerValOn)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)
		}

		{
			req := httptest.NewRequest(velocity.MethodPost, "/", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)

			req.Header.Set(headerName, headerValOff)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)

			req.Header.Set(headerName, headerValOn)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)
		}
	}

	untrustedRun := func(t *testing.T, app *velocity.App) {
		t.Helper()

		{
			req := httptest.NewRequest(velocity.MethodGet, "/", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)

			req.Header.Set(headerName, headerValOff)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)

			req.Header.Set(headerName, headerValOn)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)
		}

		{
			req := httptest.NewRequest(velocity.MethodPost, "/", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)

			req.Header.Set(headerName, headerValOff)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)

			req.Header.Set(headerName, headerValOn)
			resp, err = app.Test(req)
			require.NoError(t, err)
			require.Equal(t, velocity.StatusTooEarly, resp.StatusCode)
		}
	}

	t.Run("empty config", func(t *testing.T) {
		app := appWithConfig(t, nil)
		trustedRun(t, app)
	})
	t.Run("default config", func(t *testing.T) {
		app := appWithConfig(t, &velocity.Config{})
		trustedRun(t, app)
	})

	t.Run("config with TrustProxy", func(t *testing.T) {
		app := appWithConfig(t, &velocity.Config{
			TrustProxy: true,
		})
		untrustedRun(t, app)
	})
	t.Run("config with TrustProxy and trusted TrustProxyConfig.Proxies", func(t *testing.T) {
		app := appWithConfig(t, &velocity.Config{
			TrustProxy: true,
			TrustProxyConfig: velocity.TrustProxyConfig{
				Proxies: []string{
					"0.0.0.0",
				},
			},
		})
		trustedRun(t, app)
	})
}
