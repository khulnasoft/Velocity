package static

import (
	"embed"
	"io"
	"io/fs"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/khulnasoft/velocity"
	"github.com/stretchr/testify/require"
)

var testConfig = velocity.TestConfig{
	Timeout:       10 * time.Second,
	FailOnTimeout: true,
}

// go test -run Test_Static_Index_Default
func Test_Static_Index_Default(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/prefix", New("../../.github/workflows"))

	app.Get("", New("../../.github/"))

	app.Get("test", New("", Config{
		IndexNames: []string{"index.html"},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Hello, World!")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/not-found", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Cannot GET /not-found", string(body))
}

// go test -run Test_Static_Index
func Test_Static_Direct(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/*", New("../../.github"))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/index.html", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Hello, World!")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodPost, "/index.html", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 405, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/testdata/testRoutes.json", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMEApplicationJSON, resp.Header.Get("Content-Type"))
	require.Equal(t, "", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "test_routes")
}

// go test -run Test_Static_MaxAge
func Test_Static_MaxAge(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/*", New("../../.github", Config{
		MaxAge: 100,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/index.html", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, "text/html; charset=utf-8", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "public, max-age=100", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

// go test -run Test_Static_Custom_CacheControl
func Test_Static_Custom_CacheControl(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/*", New("../../.github", Config{
		ModifyResponse: func(c velocity.Ctx) error {
			if strings.Contains(c.GetRespHeader("Content-Type"), "text/html") {
				c.Response().Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			}
			return nil
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/index.html", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, "no-cache, no-store, must-revalidate", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")

	normalResp, normalErr := app.Test(httptest.NewRequest(velocity.MethodGet, "/config.yml", nil))
	require.NoError(t, normalErr, "app.Test(req)")
	require.Equal(t, "", normalResp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")
}

func Test_Static_Disable_Cache(t *testing.T) {
	// Skip on Windows. It's not possible to delete a file that is in use.
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	t.Parallel()

	app := velocity.New()

	file, err := os.Create("../../.github/test.txt")
	require.NoError(t, err)
	_, err = file.WriteString("Hello, World!")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	// Remove the file even if the test fails
	defer func() {
		_ = os.Remove("../../.github/test.txt") //nolint:errcheck // not needed
	}()

	app.Get("/*", New("../../.github/", Config{
		CacheDuration: -1,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/test.txt", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, "", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Hello, World!")

	require.NoError(t, os.Remove("../../.github/test.txt"))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/test.txt", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, "", resp.Header.Get(velocity.HeaderCacheControl), "CacheControl Control")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Cannot GET /test.txt", string(body))
}

func Test_Static_NotFoundHandler(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/*", New("../../.github", Config{
		NotFoundHandler: func(c velocity.Ctx) error {
			return c.SendString("Custom 404")
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/not-found", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Custom 404", string(body))
}

// go test -run Test_Static_Download
func Test_Static_Download(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/velocity.png", New("../../.github/testdata/fs/img/velocity.png", Config{
		Download: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/velocity.png", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, "image/png", resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, `attachment`, resp.Header.Get(velocity.HeaderContentDisposition))
}

// go test -run Test_Static_Group
func Test_Static_Group(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	grp := app.Group("/v1", func(c velocity.Ctx) error {
		c.Set("Test-Header", "123")
		return c.Next()
	})

	grp.Get("/v2*", New("../../.github/index.html"))

	req := httptest.NewRequest(velocity.MethodGet, "/v1/v2", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))
	require.Equal(t, "123", resp.Header.Get("Test-Header"))

	grp = app.Group("/v2")
	grp.Get("/v3*", New("../../.github/index.html"))

	req = httptest.NewRequest(velocity.MethodGet, "/v2/v3/john/doe", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))
}

func Test_Static_Wildcard(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("*", New("../../.github/index.html"))

	req := httptest.NewRequest(velocity.MethodGet, "/yesyes/john/doe", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Test file")
}

func Test_Static_Prefix_Wildcard(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/test*", New("../../.github/index.html"))

	req := httptest.NewRequest(velocity.MethodGet, "/test/john/doe", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Get("/my/nameisjohn*", New("../../.github/index.html"))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/my/nameisjohn/no/its/not", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Test file")
}

func Test_Static_Prefix(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Get("/john*", New("../../.github"))

	req := httptest.NewRequest(velocity.MethodGet, "/john/index.html", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Get("/prefix*", New("../../.github/testdata"))

	req = httptest.NewRequest(velocity.MethodGet, "/prefix/index.html", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Get("/single*", New("../../.github/testdata/testRoutes.json"))

	req = httptest.NewRequest(velocity.MethodGet, "/single", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMEApplicationJSON, resp.Header.Get(velocity.HeaderContentType))
}

func Test_Static_Trailing_Slash(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Get("/john*", New("../../.github"))

	req := httptest.NewRequest(velocity.MethodGet, "/john/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Get("/john_without_index*", New("../../.github/testdata/fs/css"))

	req = httptest.NewRequest(velocity.MethodGet, "/john_without_index/", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Use("/john", New("../../.github"))

	req = httptest.NewRequest(velocity.MethodGet, "/john/", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	req = httptest.NewRequest(velocity.MethodGet, "/john", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	app.Use("/john_without_index/", New("../../.github/testdata/fs/css"))

	req = httptest.NewRequest(velocity.MethodGet, "/john_without_index/", nil)
	resp, err = app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))
}

func Test_Static_Next(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/*", New("../../.github", Config{
		Next: func(c velocity.Ctx) bool {
			return c.Get("X-Custom-Header") == "skip"
		},
	}))

	app.Get("/*", func(c velocity.Ctx) error {
		return c.SendString("You've skipped app.Static")
	})

	t.Run("app.Static is skipped: invoking Get handler", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(velocity.MethodGet, "/", nil)
		req.Header.Set("X-Custom-Header", "skip")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
		require.Equal(t, velocity.MIMETextPlainCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "You've skipped app.Static")
	})

	t.Run("app.Static is not skipped: serving index.html", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(velocity.MethodGet, "/", nil)
		req.Header.Set("X-Custom-Header", "don't skip")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
		require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "Hello, World!")
	})
}

func Test_Route_Static_Root(t *testing.T) {
	t.Parallel()

	dir := "../../.github/testdata/fs/css"
	app := velocity.New()
	app.Get("/*", New(dir, Config{
		Browse: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")

	app = velocity.New()
	app.Get("/*", New(dir))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")
}

func Test_Route_Static_HasPrefix(t *testing.T) {
	t.Parallel()

	dir := "../../.github/testdata/fs/css"
	app := velocity.New()
	app.Get("/static*", New(dir, Config{
		Browse: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/static", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")

	app = velocity.New()
	app.Get("/static/*", New(dir, Config{
		Browse: true,
	}))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")

	app = velocity.New()
	app.Get("/static*", New(dir))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")

	app = velocity.New()
	app.Get("/static*", New(dir))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 404, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/static/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")
}

func Test_Static_FS(t *testing.T) {
	t.Parallel()

	app := velocity.New()
	app.Get("/*", New("", Config{
		FS:     os.DirFS("../../.github/testdata/fs"),
		Browse: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/css/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextCSSCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")
}

/*func Test_Static_FS_DifferentRoot(t *testing.T) {
	t.Parallel()

	app := velocity.New()
	app.Get("/*", New("fs", Config{
		FS:         os.DirFS("../../.github/testdata"),
		IndexNames: []string{"index2.html"},
		Browse:     true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "<h1>Hello, World!</h1>")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/css/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextCSSCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")
}*/

//go:embed static.go config.go
var fsTestFilesystem embed.FS

func Test_Static_FS_Browse(t *testing.T) {
	t.Parallel()

	app := velocity.New()

	app.Get("/embed*", New("", Config{
		FS:     fsTestFilesystem,
		Browse: true,
	}))

	app.Get("/dirfs*", New("", Config{
		FS:     os.DirFS("../../.github/testdata/fs/css"),
		Browse: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/dirfs", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "style.css")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/dirfs/style.css", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextCSSCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "color")

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/embed", nil))
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "app.Test(req)")
	require.Contains(t, string(body), "static.go")
}

func Test_Static_FS_Prefix_Wildcard(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Get("/test*", New("index.html", Config{
		FS:         os.DirFS("../../.github"),
		IndexNames: []string{"not_index.html"},
	}))

	req := httptest.NewRequest(velocity.MethodGet, "/test/john/doe", nil)
	resp, err := app.Test(req)
	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.NotEmpty(t, resp.Header.Get(velocity.HeaderContentLength))
	require.Equal(t, velocity.MIMETextHTMLCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Test file")
}

func Test_isFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		filesystem fs.FS
		gotError   error
		name       string
		path       string
		expected   bool
	}{
		{
			name:       "file",
			path:       "index.html",
			filesystem: os.DirFS("../../.github"),
			expected:   true,
		},
		{
			name:       "file",
			path:       "index2.html",
			filesystem: os.DirFS("../../.github"),
			expected:   false,
			gotError:   fs.ErrNotExist,
		},
		{
			name:       "directory",
			path:       ".",
			filesystem: os.DirFS("../../.github"),
			expected:   false,
		},
		{
			name:       "directory",
			path:       "not_exists",
			filesystem: os.DirFS("../../.github"),
			expected:   false,
			gotError:   fs.ErrNotExist,
		},
		{
			name:       "directory",
			path:       ".",
			filesystem: os.DirFS("../../.github/testdata/fs/css"),
			expected:   false,
		},
		{
			name:       "file",
			path:       "../../.github/testdata/fs/css/style.css",
			filesystem: nil,
			expected:   true,
		},
		{
			name:       "file",
			path:       "../../.github/testdata/fs/css/style2.css",
			filesystem: nil,
			expected:   false,
			gotError:   fs.ErrNotExist,
		},
		{
			name:       "directory",
			path:       "../../.github/testdata/fs/css",
			filesystem: nil,
			expected:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c := c
			t.Parallel()

			actual, err := isFile(c.path, c.filesystem)
			require.ErrorIs(t, err, c.gotError)
			require.Equal(t, c.expected, actual)
		})
	}
}

func Test_Static_Compress(t *testing.T) {
	t.Parallel()
	dir := "../../.github/testdata/fs" //nolint:goconst // test
	app := velocity.New()
	app.Get("/*", New(dir, Config{
		Compress: true,
	}))

	// Note: deflate is not supported by fasthttp.FS
	algorithms := []string{"zstd", "gzip", "br"}

	for _, algo := range algorithms {
		t.Run(algo+"_compression", func(t *testing.T) {
			t.Parallel()
			// request non-compressable file (less than 200 bytes), Content Lengh will remain the same
			req := httptest.NewRequest(velocity.MethodGet, "/css/style.css", nil)
			req.Header.Set("Accept-Encoding", algo)
			resp, err := app.Test(req, testConfig)

			require.NoError(t, err, "app.Test(req)")
			require.Equal(t, 200, resp.StatusCode, "Status code")
			require.Equal(t, "", resp.Header.Get(velocity.HeaderContentEncoding))
			require.Equal(t, "46", resp.Header.Get(velocity.HeaderContentLength))

			// request compressable file, ContentLenght will change
			req = httptest.NewRequest(velocity.MethodGet, "/index.html", nil)
			req.Header.Set("Accept-Encoding", algo)
			resp, err = app.Test(req, testConfig)

			require.NoError(t, err, "app.Test(req)")
			require.Equal(t, 200, resp.StatusCode, "Status code")
			require.Equal(t, algo, resp.Header.Get(velocity.HeaderContentEncoding))
			require.Greater(t, "299", resp.Header.Get(velocity.HeaderContentLength))
		})
	}
}

func Test_Static_Compress_WithoutEncoding(t *testing.T) {
	t.Parallel()
	dir := "../../.github/testdata/fs"
	app := velocity.New()
	app.Get("/*", New(dir, Config{
		Compress:      true,
		CacheDuration: 1 * time.Second,
	}))

	// request compressable file without encoding
	req := httptest.NewRequest(velocity.MethodGet, "/index.html", nil)
	resp, err := app.Test(req, testConfig)

	require.NoError(t, err, "app.Test(req)")
	require.Equal(t, 200, resp.StatusCode, "Status code")
	require.Equal(t, "", resp.Header.Get(velocity.HeaderContentEncoding))
	require.Equal(t, "299", resp.Header.Get(velocity.HeaderContentLength))

	// request compressable file with different encodings
	algorithms := []string{"zstd", "gzip", "br"}
	fileSuffixes := map[string]string{
		"gzip": ".velocity.gz",
		"br":   ".velocity.br",
		"zstd": ".velocity.zst",
	}

	for _, algo := range algorithms {
		// Wait for cache to expire
		time.Sleep(2 * time.Second)
		fileName := "index.html"
		compressedFileName := dir + "/index.html" + fileSuffixes[algo]

		req = httptest.NewRequest(velocity.MethodGet, "/"+fileName, nil)
		req.Header.Set("Accept-Encoding", algo)
		resp, err = app.Test(req, testConfig)

		require.NoError(t, err, "app.Test(req)")
		require.Equal(t, 200, resp.StatusCode, "Status code")
		require.Equal(t, algo, resp.Header.Get(velocity.HeaderContentEncoding))
		require.Greater(t, "299", resp.Header.Get(velocity.HeaderContentLength))

		// verify suffixed file was created
		_, err := os.Stat(compressedFileName)
		require.NoError(t, err, "File should exist")
	}
}

func Test_Static_Compress_WithFileSuffixes(t *testing.T) {
	t.Parallel()
	dir := "../../.github/testdata/fs"
	fileSuffixes := map[string]string{
		"gzip": ".test.gz",
		"br":   ".test.br",
		"zstd": ".test.zst",
	}

	app := velocity.New(velocity.Config{
		CompressedFileSuffixes: fileSuffixes,
	})
	app.Get("/*", New(dir, Config{
		Compress:      true,
		CacheDuration: 1 * time.Second,
	}))

	// request compressable file with different encodings
	algorithms := []string{"zstd", "gzip", "br"}

	for _, algo := range algorithms {
		// Wait for cache to expire
		time.Sleep(2 * time.Second)
		fileName := "index.html"
		compressedFileName := dir + "/index.html" + fileSuffixes[algo]

		req := httptest.NewRequest(velocity.MethodGet, "/"+fileName, nil)
		req.Header.Set("Accept-Encoding", algo)
		resp, err := app.Test(req, testConfig)

		require.NoError(t, err, "app.Test(req)")
		require.Equal(t, 200, resp.StatusCode, "Status code")
		require.Equal(t, algo, resp.Header.Get(velocity.HeaderContentEncoding))
		require.Greater(t, "299", resp.Header.Get(velocity.HeaderContentLength))

		// verify suffixed file was created
		_, err = os.Stat(compressedFileName)
		require.NoError(t, err, "File should exist")
	}
}
