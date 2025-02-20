package skip_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/skip"
)

// go test -run Test_Skip
func Test_Skip(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(skip.New(errTeapotHandler, func(velocity.Ctx) bool { return true }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
}

// go test -run Test_SkipFalse
func Test_SkipFalse(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(skip.New(errTeapotHandler, func(velocity.Ctx) bool { return false }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusTeapot, resp.StatusCode)
}

// go test -run Test_SkipNilFunc
func Test_SkipNilFunc(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(skip.New(errTeapotHandler, nil))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusTeapot, resp.StatusCode)
}

func helloWorldHandler(c velocity.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func errTeapotHandler(velocity.Ctx) error {
	return velocity.ErrTeapot
}
