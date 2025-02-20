package expvar

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity"
)

func Test_Non_Expvar_Path(t *testing.T) {
	t.Parallel()
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

func Test_Expvar_Index(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/vars", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, velocity.MIMEApplicationJSONCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.True(t, bytes.Contains(b, []byte("cmdline")))
	require.True(t, bytes.Contains(b, []byte("memstat")))
}

func Test_Expvar_Filter(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/vars?r=cmd", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, velocity.MIMEApplicationJSONCharsetUTF8, resp.Header.Get(velocity.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.True(t, bytes.Contains(b, []byte("cmdline")))
	require.False(t, bytes.Contains(b, []byte("memstat")))
}

func Test_Expvar_Other_Path(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/vars/302", nil))
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
}

// go test -run Test_Expvar_Next
func Test_Expvar_Next(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/debug/vars", nil))
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)
}
