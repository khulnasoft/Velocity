package requestid

import (
	"net/http/httptest"
	"testing"

	"github.com/khulnasoft/velocity"
	"github.com/stretchr/testify/require"
)

// go test -run Test_RequestID
func Test_RequestID(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)

	reqid := resp.Header.Get(velocity.HeaderXRequestID)
	require.Len(t, reqid, 36)

	req := httptest.NewRequest(velocity.MethodGet, "/", nil)
	req.Header.Add(velocity.HeaderXRequestID, reqid)

	resp, err = app.Test(req)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, reqid, resp.Header.Get(velocity.HeaderXRequestID))
}

// go test -run Test_RequestID_Next
func Test_RequestID_Next(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Empty(t, resp.Header.Get(velocity.HeaderXRequestID))
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
}

// go test -run Test_RequestID_FromContext
func Test_RequestID_FromContext(t *testing.T) {
	t.Parallel()

	reqID := "ThisIsARequestId"

	type args struct {
		inputFunc func(c velocity.Ctx) any
	}

	tests := []struct {
		args args
		name string
	}{
		{
			name: "From velocity.Ctx",
			args: args{
				inputFunc: func(c velocity.Ctx) any {
					return c
				},
			},
		},
		{
			name: "From context.Context",
			args: args{
				inputFunc: func(c velocity.Ctx) any {
					return c.Context()
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := velocity.New()
			app.Use(New(Config{
				Generator: func() string {
					return reqID
				},
			}))

			var ctxVal string

			app.Use(func(c velocity.Ctx) error {
				ctxVal = FromContext(tt.args.inputFunc(c))
				return c.Next()
			})

			_, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
			require.NoError(t, err)
			require.Equal(t, reqID, ctxVal)
		})
	}
}
