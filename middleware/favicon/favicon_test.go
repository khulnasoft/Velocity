package favicon

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/khulnasoft/velocity"
)

// go test -run Test_Middleware_Favicon
func Test_Middleware_Favicon(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(_ velocity.Ctx) error {
		return nil
	})

	// Skip Favicon middleware
	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusNoContent, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodOptions, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodPut, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusMethodNotAllowed, resp.StatusCode, "Status code")
	require.Equal(t, "GET, HEAD, OPTIONS", resp.Header.Get(velocity.HeaderAllow))
}

// go test -run Test_Middleware_Favicon_Not_Found
func Test_Middleware_Favicon_Not_Found(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Error("should cache panic")
			return
		}
	}()

	velocity.New().Use(New(Config{
		File: "non-exist.ico",
	}))
}

// go test -run Test_Middleware_Favicon_Found
func Test_Middleware_Favicon_Found(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		File: "../../.github/testdata/favicon.ico",
	}))

	app.Get("/", func(_ velocity.Ctx) error {
		return nil
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")
	require.Equal(t, "image/x-icon", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "public, max-age=31536000", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

// go test -run Test_Custom_Favicon_Url
func Test_Custom_Favicon_URL(t *testing.T) {
	app := velocity.New()
	const customURL = "/favicon.svg"
	app.Use(New(Config{
		File: "../../.github/testdata/favicon.ico",
		URL:  customURL,
	}))

	app.Get("/", func(_ velocity.Ctx) error {
		return nil
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, customURL, nil))

	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")
	require.Equal(t, "image/x-icon", resp.Header.Get(velocity.HeaderContentType))
}

// go test -run Test_Custom_Favicon_Data
func Test_Custom_Favicon_Data(t *testing.T) {
	data, err := os.ReadFile("../../.github/testdata/favicon.ico")
	require.NoError(t, err)

	app := velocity.New()

	app.Use(New(Config{
		Data: data,
	}))

	app.Get("/", func(_ velocity.Ctx) error {
		return nil
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")
	require.Equal(t, "image/x-icon", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "public, max-age=31536000", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

// go test -run Test_Middleware_Favicon_FileSystem
func Test_Middleware_Favicon_FileSystem(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		File:       "favicon.ico",
		FileSystem: os.DirFS("../../.github/testdata"),
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")
	require.Equal(t, "image/x-icon", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "public, max-age=31536000", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

// go test -run Test_Middleware_Favicon_CacheControl
func Test_Middleware_Favicon_CacheControl(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		CacheControl: "public, max-age=100",
		File:         "../../.github/testdata/favicon.ico",
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/favicon.ico", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, velocity.StatusOK, resp.StatusCode, "Status code")
	require.Equal(t, "image/x-icon", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "public, max-age=100", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

// go test -v -run=^$ -bench=Benchmark_Middleware_Favicon -benchmem -count=4
func Benchmark_Middleware_Favicon(b *testing.B) {
	app := velocity.New()
	app.Use(New())
	app.Get("/", func(_ velocity.Ctx) error {
		return nil
	})
	handler := app.Handler()

	c := &fasthttp.RequestCtx{}
	c.Request.SetRequestURI("/")

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler(c)
	}
}

// go test -run Test_Favicon_Next
func Test_Favicon_Next(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
}
