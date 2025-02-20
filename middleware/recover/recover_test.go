package recover //nolint:predeclared // TODO: Rename to some non-builtin

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity"
)

// go test -run Test_Recover
func Test_Recover(t *testing.T) {
	t.Parallel()
	app := velocity.New(velocity.Config{
		ErrorHandler: func(c velocity.Ctx, err error) error {
			require.Equal(t, "Hi, I'm an error!", err.Error())
			return c.SendStatus(velocity.StatusTeapot)
		},
	})

	app.Use(New())

	app.Get("/panic", func(_ velocity.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/panic", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusTeapot, resp.StatusCode)
}

// go test -run Test_Recover_Next
func Test_Recover_Next(t *testing.T) {
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

func Test_Recover_EnableStackTrace(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Use(New(Config{
		EnableStackTrace: true,
	}))

	app.Get("/panic", func(_ velocity.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/panic", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusInternalServerError, resp.StatusCode)
}
