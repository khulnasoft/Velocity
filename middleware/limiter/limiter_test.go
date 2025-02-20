package limiter

import (
	"io"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/internal/storage/memory"
)

// go test -run Test_Limiter_With_Max_Func_With_Zero -race -v
func Test_Limiter_With_Max_Func_With_Zero_And_Limiter_Sliding(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		MaxFunc:                func(_ velocity.Ctx) int { return 0 },
		Expiration:             2 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" { //nolint:goconst // test
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_With_Max_Func_With_Zero -race -v
func Test_Limiter_With_Max_Func_With_Zero(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		MaxFunc: func(_ velocity.Ctx) int {
			return 0
		},
		Expiration: 2 * time.Second,
		Storage:    memory.New(),
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	var wg sync.WaitGroup

	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
			assert.NoError(t, err)
			assert.Equal(t, velocity.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Hello tester!", string(body))
		}(&wg)
	}

	wg.Wait()

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_With_Max_Func -race -v
func Test_Limiter_With_Max_Func(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	maxRequests := 10

	app.Use(New(Config{
		MaxFunc: func(_ velocity.Ctx) int {
			return maxRequests
		},
		Expiration: 2 * time.Second,
		Storage:    memory.New(),
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	var wg sync.WaitGroup

	for i := 0; i <= maxRequests-1; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
			assert.NoError(t, err)
			assert.Equal(t, velocity.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Hello tester!", string(body))
		}(&wg)
	}

	wg.Wait()

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Concurrency_Store -race -v
func Test_Limiter_Concurrency_Store(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:        50,
		Expiration: 2 * time.Second,
		Storage:    memory.New(),
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	var wg sync.WaitGroup

	for i := 0; i <= 49; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
			assert.NoError(t, err)
			assert.Equal(t, velocity.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Hello tester!", string(body))
		}(&wg)
	}

	wg.Wait()

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Concurrency -race -v
func Test_Limiter_Concurrency(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:        50,
		Expiration: 2 * time.Second,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	var wg sync.WaitGroup

	for i := 0; i <= 49; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
			assert.NoError(t, err)
			assert.Equal(t, velocity.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Hello tester!", string(body))
		}(&wg)
	}

	wg.Wait()

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_No_Skip_Choices -v
func Test_Limiter_Fixed_Window_No_Skip_Choices(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    2,
		Expiration:             2 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_Custom_Storage_No_Skip_Choices -v
func Test_Limiter_Fixed_Window_Custom_Storage_No_Skip_Choices(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    2,
		Expiration:             2 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                memory.New(),
		LimiterMiddleware:      FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_No_Skip_Choices -v
func Test_Limiter_Sliding_Window_No_Skip_Choices(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    2,
		Expiration:             2 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_Custom_Storage_No_Skip_Choices -v
func Test_Limiter_Sliding_Window_Custom_Storage_No_Skip_Choices(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    2,
		Expiration:             2 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                memory.New(),
		LimiterMiddleware:      SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_Skip_Failed_Requests -v
func Test_Limiter_Fixed_Window_Skip_Failed_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                1,
		Expiration:         2 * time.Second,
		SkipFailedRequests: true,
		LimiterMiddleware:  FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_Custom_Storage_Skip_Failed_Requests -v
func Test_Limiter_Fixed_Window_Custom_Storage_Skip_Failed_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                1,
		Expiration:         2 * time.Second,
		Storage:            memory.New(),
		SkipFailedRequests: true,
		LimiterMiddleware:  FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_Skip_Failed_Requests -v
func Test_Limiter_Sliding_Window_Skip_Failed_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                1,
		Expiration:         2 * time.Second,
		SkipFailedRequests: true,
		LimiterMiddleware:  SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_Custom_Storage_Skip_Failed_Requests -v
func Test_Limiter_Sliding_Window_Custom_Storage_Skip_Failed_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                1,
		Expiration:         2 * time.Second,
		Storage:            memory.New(),
		SkipFailedRequests: true,
		LimiterMiddleware:  SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_Skip_Successful_Requests -v
func Test_Limiter_Fixed_Window_Skip_Successful_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    1,
		Expiration:             2 * time.Second,
		SkipSuccessfulRequests: true,
		LimiterMiddleware:      FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// go test -run Test_Limiter_Fixed_Window_Custom_Storage_Skip_Successful_Requests -v
func Test_Limiter_Fixed_Window_Custom_Storage_Skip_Successful_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    1,
		Expiration:             2 * time.Second,
		Storage:                memory.New(),
		SkipSuccessfulRequests: true,
		LimiterMiddleware:      FixedWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(3 * time.Second)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_Skip_Successful_Requests -v
func Test_Limiter_Sliding_Window_Skip_Successful_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    1,
		Expiration:             2 * time.Second,
		SkipSuccessfulRequests: true,
		LimiterMiddleware:      SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// go test -run Test_Limiter_Sliding_Window_Custom_Storage_Skip_Successful_Requests -v
func Test_Limiter_Sliding_Window_Custom_Storage_Skip_Successful_Requests(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:                    1,
		Expiration:             2 * time.Second,
		Storage:                memory.New(),
		SkipSuccessfulRequests: true,
		LimiterMiddleware:      SlidingWindow{},
	}))

	app.Get("/:status", func(c velocity.Ctx) error {
		if c.Params("status") == "fail" {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/success", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	time.Sleep(4*time.Second + 500*time.Millisecond)

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// go test -v -run=^$ -bench=Benchmark_Limiter_Custom_Store -benchmem -count=4
func Benchmark_Limiter_Custom_Store(b *testing.B) {
	app := velocity.New()

	app.Use(New(Config{
		Max:        100,
		Expiration: 60 * time.Second,
		Storage:    memory.New(),
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(velocity.MethodGet)
	fctx.Request.SetRequestURI("/")

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h(fctx)
	}
}

// go test -run Test_Limiter_Next
func Test_Limiter_Next(t *testing.T) {
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

func Test_Limiter_Headers(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Max:        50,
		Expiration: 2 * time.Second,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(velocity.MethodGet)
	fctx.Request.SetRequestURI("/")

	app.Handler()(fctx)

	require.Equal(t, "50", string(fctx.Response.Header.Peek("X-RateLimit-Limit")))
	if v := string(fctx.Response.Header.Peek("X-RateLimit-Remaining")); v == "" {
		t.Errorf("The X-RateLimit-Remaining header is not set correctly - value is an empty string.")
	}
	if v := string(fctx.Response.Header.Peek("X-RateLimit-Reset")); !(v == "1" || v == "2") {
		t.Errorf("The X-RateLimit-Reset header is not set correctly - value is out of bounds.")
	}
}

// go test -v -run=^$ -bench=Benchmark_Limiter -benchmem -count=4
func Benchmark_Limiter(b *testing.B) {
	app := velocity.New()

	app.Use(New(Config{
		Max:        100,
		Expiration: 60 * time.Second,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(velocity.MethodGet)
	fctx.Request.SetRequestURI("/")

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h(fctx)
	}
}

// go test -run Test_Sliding_Window -race -v
func Test_Sliding_Window(t *testing.T) {
	t.Parallel()
	app := velocity.New()
	app.Use(New(Config{
		Max:               10,
		Expiration:        1 * time.Second,
		Storage:           memory.New(),
		LimiterMiddleware: SlidingWindow{},
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello tester!")
	})

	singleRequest := func(shouldFail bool) {
		resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
		if shouldFail {
			require.NoError(t, err)
			require.Equal(t, 429, resp.StatusCode)
		} else {
			require.NoError(t, err)
			require.Equal(t, velocity.StatusOK, resp.StatusCode)
		}
	}

	for i := 0; i < 5; i++ {
		singleRequest(false)
	}

	time.Sleep(3 * time.Second)

	for i := 0; i < 5; i++ {
		singleRequest(false)
	}

	time.Sleep(3 * time.Second)

	for i := 0; i < 5; i++ {
		singleRequest(false)
	}

	time.Sleep(3 * time.Second)

	for i := 0; i < 10; i++ {
		singleRequest(false)
	}

	// requests should fail now
	for i := 0; i < 5; i++ {
		singleRequest(true)
	}
}
