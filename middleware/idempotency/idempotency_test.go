package idempotency_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/idempotency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go test -run Test_Idempotency
func Test_Idempotency(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(func(c velocity.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}

		isMethodSafe := velocity.IsMethodSafe(c.Method())
		isIdempotent := idempotency.IsFromCache(c) || idempotency.WasPutToCache(c)
		hasReqHeader := c.Get("X-Idempotency-Key") != ""

		if isMethodSafe {
			if isIdempotent {
				return errors.New("request with safe HTTP method should not be idempotent")
			}
		} else {
			// Unsafe
			if hasReqHeader {
				if !isIdempotent {
					return errors.New("request with unsafe HTTP method should be idempotent if X-Idempotency-Key request header is set")
				}
			} else if isIdempotent {
				return errors.New("request with unsafe HTTP method should not be idempotent if X-Idempotency-Key request header is not set")
			}
		}

		return nil
	})

	// Needs to be at least a second as the memory storage doesn't support shorter durations.
	const lifetime = 2 * time.Second

	app.Use(idempotency.New(idempotency.Config{
		Lifetime: lifetime,
	}))

	nextCount := func() func() int {
		var count int32
		return func() int {
			return int(atomic.AddInt32(&count, 1))
		}
	}()

	app.Add([]string{
		velocity.MethodGet,
		velocity.MethodPost,
	}, "/", func(c velocity.Ctx) error {
		return c.SendString(strconv.Itoa(nextCount()))
	})

	app.Post("/slow", func(c velocity.Ctx) error {
		time.Sleep(3 * lifetime)

		return c.SendString(strconv.Itoa(nextCount()))
	})

	doReq := func(method, route, idempotencyKey string) string {
		req := httptest.NewRequest(method, route, nil)
		if idempotencyKey != "" {
			req.Header.Set("X-Idempotency-Key", idempotencyKey)
		}
		resp, err := app.Test(req, velocity.TestConfig{
			Timeout:       15 * time.Second,
			FailOnTimeout: true,
		})
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, velocity.StatusOK, resp.StatusCode, string(body))
		return string(body)
	}

	require.Equal(t, "1", doReq(velocity.MethodGet, "/", ""))
	require.Equal(t, "2", doReq(velocity.MethodGet, "/", ""))

	require.Equal(t, "3", doReq(velocity.MethodPost, "/", ""))
	require.Equal(t, "4", doReq(velocity.MethodPost, "/", ""))

	require.Equal(t, "5", doReq(velocity.MethodGet, "/", "00000000-0000-0000-0000-000000000000"))
	require.Equal(t, "6", doReq(velocity.MethodGet, "/", "00000000-0000-0000-0000-000000000000"))

	require.Equal(t, "7", doReq(velocity.MethodPost, "/", "00000000-0000-0000-0000-000000000000"))
	require.Equal(t, "7", doReq(velocity.MethodPost, "/", "00000000-0000-0000-0000-000000000000"))
	require.Equal(t, "8", doReq(velocity.MethodPost, "/", ""))
	require.Equal(t, "9", doReq(velocity.MethodPost, "/", "11111111-1111-1111-1111-111111111111"))

	require.Equal(t, "7", doReq(velocity.MethodPost, "/", "00000000-0000-0000-0000-000000000000"))
	time.Sleep(4 * lifetime)
	require.Equal(t, "10", doReq(velocity.MethodPost, "/", "00000000-0000-0000-0000-000000000000"))
	require.Equal(t, "10", doReq(velocity.MethodPost, "/", "00000000-0000-0000-0000-000000000000"))

	// Test raciness
	{
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.Equal(t, "11", doReq(velocity.MethodPost, "/slow", "22222222-2222-2222-2222-222222222222"))
			}()
		}
		wg.Wait()
		require.Equal(t, "11", doReq(velocity.MethodPost, "/slow", "22222222-2222-2222-2222-222222222222"))
	}
	time.Sleep(3 * lifetime)
	require.Equal(t, "12", doReq(velocity.MethodPost, "/slow", "22222222-2222-2222-2222-222222222222"))
}

// go test -v -run=^$ -bench=Benchmark_Idempotency -benchmem -count=4
func Benchmark_Idempotency(b *testing.B) {
	app := velocity.New()

	// Needs to be at least a second as the memory storage doesn't support shorter durations.
	const lifetime = 1 * time.Second

	app.Use(idempotency.New(idempotency.Config{
		Lifetime: lifetime,
	}))

	app.Post("/", func(_ velocity.Ctx) error {
		return nil
	})

	h := app.Handler()

	b.Run("hit", func(b *testing.B) {
		c := &fasthttp.RequestCtx{}
		c.Request.Header.SetMethod(velocity.MethodPost)
		c.Request.SetRequestURI("/")
		c.Request.Header.Set("X-Idempotency-Key", "00000000-0000-0000-0000-000000000000")

		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			h(c)
		}
	})

	b.Run("skip", func(b *testing.B) {
		c := &fasthttp.RequestCtx{}
		c.Request.Header.SetMethod(velocity.MethodPost)
		c.Request.SetRequestURI("/")

		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			h(c)
		}
	})
}
