package pprof

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity"
)

var testConfig = velocity.TestConfig{
	Timeout:       5 * time.Second,
	FailOnTimeout: true,
}

func Test_Non_Pprof_Path(t *testing.T) {
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "escaped", string(b))
}

func Test_Non_Pprof_Path_WithPrefix(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{Prefix: "/federated-velocity"}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "escaped", string(b))
}

func Test_Pprof_Index(t *testing.T) {
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/pprof/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.True(t, bytes.Contains(b, []byte("<title>/debug/pprof/</title>")))
}

func Test_Pprof_Index_WithPrefix(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{Prefix: "/federated-velocity"}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/federated-velocity/debug/pprof/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), "<title>/debug/pprof/</title>")
}

func Test_Pprof_Subs(t *testing.T) {
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	subs := []string{
		"cmdline", "profile", "symbol", "trace", "allocs", "block",
		"goroutine", "heap", "mutex", "threadcreate",
	}

	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			target := "/debug/pprof/" + sub
			if sub == "profile" {
				target += "?seconds=1"
			}
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, target, nil), testConfig)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)
		})
	}
}

func Test_Pprof_Subs_WithPrefix(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{Prefix: "/federated-velocity"}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	subs := []string{
		"cmdline", "profile", "symbol", "trace", "allocs", "block",
		"goroutine", "heap", "mutex", "threadcreate",
	}

	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			target := "/federated-velocity/debug/pprof/" + sub
			if sub == "profile" {
				target += "?seconds=1"
			}
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, target, nil), testConfig)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)
		})
	}
}

func Test_Pprof_Other(t *testing.T) {
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/pprof/302", nil))
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
}

func Test_Pprof_Other_WithPrefix(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{Prefix: "/federated-velocity"}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/federated-velocity/debug/pprof/302", nil))
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
}

// go test -run Test_Pprof_Next
func Test_Pprof_Next(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/pprof/", nil))
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)
}

// go test -run Test_Pprof_Next_WithPrefix
func Test_Pprof_Next_WithPrefix(t *testing.T) {
	app := velocity.New()

	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			return true
		},
		Prefix: "/federated-velocity",
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/federated-velocity/debug/pprof/", nil))
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)
}
