//nolint:depguard // Because we test logging :D
package logger

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
	velocitylog "go.khulnasoft.com/velocity/log"
	"go.khulnasoft.com/velocity/middleware/requestid"
)

func benchmarkSetup(b *testing.B, app *velocity.App, uri string) {
	b.Helper()

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(velocity.MethodGet)
	fctx.Request.SetRequestURI(uri)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h(fctx)
	}
}

func benchmarkSetupParallel(b *testing.B, app *velocity.App, path string) {
	b.Helper()

	handler := app.Handler()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		fctx := &fasthttp.RequestCtx{}
		fctx.Request.Header.SetMethod(velocity.MethodGet)
		fctx.Request.SetRequestURI(path)

		for pb.Next() {
			handler(fctx)
		}
	})
}

// go test -run Test_Logger
func Test_Logger(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Use(New(Config{
		Format: "${error}",
		Output: buf,
	}))

	app.Get("/", func(_ velocity.Ctx) error {
		return errors.New("some random error")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, "some random error", buf.String())
}

// go test -run Test_Logger_locals
func Test_Logger_locals(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Use(New(Config{
		Format: "${locals:demo}",
		Output: buf,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		c.Locals("demo", "johndoe")
		return c.SendStatus(velocity.StatusOK)
	})

	app.Get("/int", func(c velocity.Ctx) error {
		c.Locals("demo", 55)
		return c.SendStatus(velocity.StatusOK)
	})

	app.Get("/empty", func(c velocity.Ctx) error {
		return c.SendStatus(velocity.StatusOK)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "johndoe", buf.String())

	buf.Reset()

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/int", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "55", buf.String())

	buf.Reset()

	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/empty", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "", buf.String())
}

// go test -run Test_Logger_Next
func Test_Logger_Next(t *testing.T) {
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

// go test -run Test_Logger_Done
func Test_Logger_Done(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer(nil)
	app := velocity.New()

	app.Use(New(Config{
		Done: func(c velocity.Ctx, logString []byte) {
			if c.Response().StatusCode() == velocity.StatusOK {
				_, err := buf.Write(logString)
				require.NoError(t, err)
			}
		},
	})).Get("/logging", func(ctx velocity.Ctx) error {
		return ctx.SendStatus(velocity.StatusOK)
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/logging", nil))

	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Positive(t, buf.Len(), 0)
}

// go test -run Test_Logger_ErrorTimeZone
func Test_Logger_ErrorTimeZone(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		TimeZone: "invalid",
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
}

// go test -run Test_Logger_Velocity_Logger
func Test_Logger_LoggerToWriter(t *testing.T) {
	app := velocity.New()

	buf := bytebufferpool.Get()
	t.Cleanup(func() {
		bytebufferpool.Put(buf)
	})

	logger := velocitylog.DefaultLogger()
	stdlogger, ok := logger.Logger().(*log.Logger)
	require.True(t, ok)

	stdlogger.SetFlags(0)
	logger.SetOutput(buf)

	testCases := []struct {
		levelStr string
		level    velocitylog.Level
	}{
		{
			level:    velocitylog.LevelTrace,
			levelStr: "Trace",
		},
		{
			level:    velocitylog.LevelDebug,
			levelStr: "Debug",
		},
		{
			level:    velocitylog.LevelInfo,
			levelStr: "Info",
		},
		{
			level:    velocitylog.LevelWarn,
			levelStr: "Warn",
		},
		{
			level:    velocitylog.LevelError,
			levelStr: "Error",
		},
	}

	for _, tc := range testCases {
		level := strconv.Itoa(int(tc.level))
		t.Run(level, func(t *testing.T) {
			buf.Reset()

			app.Use("/"+level, New(Config{
				Format: "${error}",
				Output: LoggerToWriter(logger, tc.
					level),
			}))

			app.Get("/"+level, func(_ velocity.Ctx) error {
				return errors.New("some random error")
			})

			resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/"+level, nil))
			require.NoError(t, err)
			require.Equal(t, velocity.StatusInternalServerError, resp.StatusCode)
			require.Equal(t, "["+tc.levelStr+"] some random error\n", buf.String())
		})

		require.Panics(t, func() {
			LoggerToWriter(logger, velocitylog.LevelPanic)
		})

		require.Panics(t, func() {
			LoggerToWriter(logger, velocitylog.LevelFatal)
		})

		require.Panics(t, func() {
			LoggerToWriter(nil, velocitylog.LevelFatal)
		})
	}
}

type fakeErrorOutput int

func (o *fakeErrorOutput) Write([]byte) (int, error) {
	*o++
	return 0, errors.New("fake output")
}

// go test -run Test_Logger_ErrorOutput_WithoutColor
func Test_Logger_ErrorOutput_WithoutColor(t *testing.T) {
	t.Parallel()
	o := new(fakeErrorOutput)
	app := velocity.New()

	app.Use(New(Config{
		Output:        o,
		DisableColors: true,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
	require.EqualValues(t, 2, *o)
}

// go test -run Test_Logger_ErrorOutput
func Test_Logger_ErrorOutput(t *testing.T) {
	t.Parallel()
	o := new(fakeErrorOutput)
	app := velocity.New()

	app.Use(New(Config{
		Output: o,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
	require.EqualValues(t, 2, *o)
}

// go test -run Test_Logger_All
func Test_Logger_All(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${pid}${reqHeaders}${referer}${scheme}${protocol}${ip}${ips}${host}${url}${ua}${body}${route}${black}${red}${green}${yellow}${blue}${magenta}${cyan}${white}${reset}${error}${reqHeader:test}${query:test}${form:test}${cookie:test}${non}",
		Output: buf,
	}))

	// Alias colors
	colors := app.Config().ColorScheme

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/?foo=bar", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)

	expected := fmt.Sprintf("%dHost=example.comhttpHTTP/1.10.0.0.0example.com/?foo=bar/%s%s%s%s%s%s%s%s%sCannot GET /", os.Getpid(), colors.Black, colors.Red, colors.Green, colors.Yellow, colors.Blue, colors.Magenta, colors.Cyan, colors.White, colors.Reset)
	require.Equal(t, expected, buf.String())
}

func getLatencyTimeUnits() []struct {
	unit string
	div  time.Duration
} {
	// windows does not support µs sleep precision
	// https://github.com/golang/go/issues/29485
	if runtime.GOOS == "windows" {
		return []struct {
			unit string
			div  time.Duration
		}{
			{unit: "ms", div: time.Millisecond},
			{unit: "s", div: time.Second},
		}
	}
	return []struct {
		unit string
		div  time.Duration
	}{
		{unit: "µs", div: time.Microsecond},
		{unit: "ms", div: time.Millisecond},
		{unit: "s", div: time.Second},
	}
}

// go test -run Test_Logger_WithLatency
func Test_Logger_WithLatency(t *testing.T) {
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)
	app := velocity.New()

	logger := New(Config{
		Output: buff,
		Format: "${latency}",
	})
	app.Use(logger)

	// Define a list of time units to test
	timeUnits := getLatencyTimeUnits()

	// Initialize a new time unit
	sleepDuration := 1 * time.Nanosecond

	// Define a test route that sleeps
	app.Get("/test", func(c velocity.Ctx) error {
		time.Sleep(sleepDuration)
		return c.SendStatus(velocity.StatusOK)
	})

	// Loop through each time unit and assert that the log output contains the expected latency value
	for _, tu := range timeUnits {
		// Update the sleep duration for the next iteration
		sleepDuration = 1 * tu.div

		// Create a new HTTP request to the test route
		resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/test", nil), velocity.TestConfig{
			Timeout:       3 * time.Second,
			FailOnTimeout: true,
		})
		require.NoError(t, err)
		require.Equal(t, velocity.StatusOK, resp.StatusCode)

		// Assert that the log output contains the expected latency value in the current time unit
		require.True(t, bytes.HasSuffix(buff.Bytes(), []byte(tu.unit)), "Expected latency to be in %s, got %s", tu.unit, buff.String())

		// Reset the buffer
		buff.Reset()
	}
}

// go test -run Test_Logger_WithLatency_DefaultFormat
func Test_Logger_WithLatency_DefaultFormat(t *testing.T) {
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)
	app := velocity.New()

	logger := New(Config{
		Output: buff,
	})
	app.Use(logger)

	// Define a list of time units to test
	timeUnits := getLatencyTimeUnits()

	// Initialize a new time unit
	sleepDuration := 1 * time.Nanosecond

	// Define a test route that sleeps
	app.Get("/test", func(c velocity.Ctx) error {
		time.Sleep(sleepDuration)
		return c.SendStatus(velocity.StatusOK)
	})

	// Loop through each time unit and assert that the log output contains the expected latency value
	for _, tu := range timeUnits {
		// Update the sleep duration for the next iteration
		sleepDuration = 1 * tu.div

		// Create a new HTTP request to the test route
		resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/test", nil), velocity.TestConfig{
			Timeout:       2 * time.Second,
			FailOnTimeout: true,
		})
		require.NoError(t, err)
		require.Equal(t, velocity.StatusOK, resp.StatusCode)

		// Assert that the log output contains the expected latency value in the current time unit
		// parse out the latency value from the log output
		latency := bytes.Split(buff.Bytes(), []byte(" | "))[2]
		// Assert that the latency value is in the current time unit
		require.True(t, bytes.HasSuffix(latency, []byte(tu.unit)), "Expected latency to be in %s, got %s", tu.unit, latency)

		// Reset the buffer
		buff.Reset()
	}
}

// go test -run Test_Query_Params
func Test_Query_Params(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${queryParams}",
		Output: buf,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/?foo=bar&baz=moz", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)

	expected := "foo=bar&baz=moz"
	require.Equal(t, expected, buf.String())
}

// go test -run Test_Response_Body
func Test_Response_Body(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${resBody}",
		Output: buf,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Sample response body")
	})

	app.Post("/test", func(c velocity.Ctx) error {
		return c.Send([]byte("Post in test"))
	})

	_, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)

	expectedGetResponse := "Sample response body"
	require.Equal(t, expectedGetResponse, buf.String())

	buf.Reset() // Reset buffer to test POST
	_, err = app.Test(httptest.NewRequest(velocity.MethodPost, "/test", nil))

	expectedPostResponse := "Post in test"
	require.NoError(t, err)
	require.Equal(t, expectedPostResponse, buf.String())
}

// go test -run Test_Request_Body
func Test_Request_Body(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	app := velocity.New()

	app.Use(New(Config{
		Format: "${bytesReceived} ${bytesSent} ${status}",
		Output: buf,
	}))

	app.Post("/", func(c velocity.Ctx) error {
		c.Response().Header.SetContentLength(5)
		return c.SendString("World")
	})

	// Create a POST request with a body
	body := []byte("Hello")
	req := httptest.NewRequest(velocity.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/octet-stream")

	_, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, "5 5 200", buf.String())
}

// go test -run Test_Logger_AppendUint
func Test_Logger_AppendUint(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Use(New(Config{
		Format: "${bytesReceived} ${bytesSent} ${status}",
		Output: buf,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("hello")
	})

	app.Get("/content", func(c velocity.Ctx) error {
		c.Response().Header.SetContentLength(5)
		return c.SendString("hello")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "-2 0 200", buf.String())

	buf.Reset()
	resp, err = app.Test(httptest.NewRequest(velocity.MethodGet, "/content", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "-2 5 200", buf.String())
}

// go test -run Test_Logger_Data_Race -race
func Test_Logger_Data_Race(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Use(New(ConfigDefault))
	app.Use(New(Config{
		Format: "${time} | ${pid} | ${locals:requestid} | ${status} | ${latency} | ${method} | ${path}\n",
	}))

	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("hello")
	})

	var (
		resp1, resp2 *http.Response
		err1, err2   error
	)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		resp1, err1 = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
		wg.Done()
	}()
	resp2, err2 = app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	wg.Wait()

	require.NoError(t, err1)
	require.Equal(t, velocity.StatusOK, resp1.StatusCode)
	require.NoError(t, err2)
	require.Equal(t, velocity.StatusOK, resp2.StatusCode)
}

// go test -run Test_Response_Header
func Test_Response_Header(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(requestid.New(requestid.Config{
		Next:      nil,
		Header:    velocity.HeaderXRequestID,
		Generator: func() string { return "Hello velocity!" },
	}))
	app.Use(New(Config{
		Format: "${respHeader:X-Request-ID}",
		Output: buf,
	}))
	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello velocity!")
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))

	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello velocity!", buf.String())
}

// go test -run Test_Req_Header
func Test_Req_Header(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${reqHeader:test}",
		Output: buf,
	}))
	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello velocity!")
	})
	headerReq := httptest.NewRequest(velocity.MethodGet, "/", nil)
	headerReq.Header.Add("test", "Hello velocity!")

	resp, err := app.Test(headerReq)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello velocity!", buf.String())
}

// go test -run Test_ReqHeader_Header
func Test_ReqHeader_Header(t *testing.T) {
	t.Parallel()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${reqHeader:test}",
		Output: buf,
	}))
	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello velocity!")
	})
	reqHeaderReq := httptest.NewRequest(velocity.MethodGet, "/", nil)
	reqHeaderReq.Header.Add("test", "Hello velocity!")

	resp, err := app.Test(reqHeaderReq)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello velocity!", buf.String())
}

// go test -run Test_CustomTags
func Test_CustomTags(t *testing.T) {
	t.Parallel()
	customTag := "it is a custom tag"

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app := velocity.New()

	app.Use(New(Config{
		Format: "${custom_tag}",
		CustomTags: map[string]LogFunc{
			"custom_tag": func(output Buffer, _ velocity.Ctx, _ *Data, _ string) (int, error) {
				return output.WriteString(customTag)
			},
		},
		Output: buf,
	}))
	app.Get("/", func(c velocity.Ctx) error {
		return c.SendString("Hello velocity!")
	})
	reqHeaderReq := httptest.NewRequest(velocity.MethodGet, "/", nil)
	reqHeaderReq.Header.Add("test", "Hello velocity!")

	resp, err := app.Test(reqHeaderReq)
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)
	require.Equal(t, customTag, buf.String())
}

// go test -run Test_Logger_ByteSent_Streaming
func Test_Logger_ByteSent_Streaming(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Use(New(Config{
		Format: "${bytesReceived} ${bytesSent} ${status}",
		Output: buf,
	}))

	app.Get("/", func(c velocity.Ctx) error {
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				fmt.Fprintf(w, "data: Message: %s\n\n", msg) //nolint:errcheck // ignore error
				err := w.Flush()
				if err != nil {
					break
				}
				if i == 10 {
					break
				}
			}
		})
		return nil
	})

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusOK, resp.StatusCode)

	// -2 means identity, -1 means chunked, 200 status
	require.Equal(t, "-2 -1 200", buf.String())
}

type fakeOutput int

func (o *fakeOutput) Write(b []byte) (int, error) {
	*o++
	return len(b), nil
}

// go test -run Test_Logger_EnableColors
func Test_Logger_EnableColors(t *testing.T) {
	t.Parallel()
	o := new(fakeOutput)
	app := velocity.New()

	app.Use(New(Config{
		Output: o,
	}))

	resp, err := app.Test(httptest.NewRequest(velocity.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, velocity.StatusNotFound, resp.StatusCode)
	require.EqualValues(t, 1, *o)
}

// go test -v -run=^$ -bench=Benchmark_Logger$ -benchmem -count=4
func Benchmark_Logger(b *testing.B) {
	b.Run("NoMiddleware", func(bb *testing.B) {
		app := velocity.New()
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("WithBytesAndStatus", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("test", "test")
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("DefaultFormat", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("DefaultFormatDisableColors", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Output:        io.Discard,
			DisableColors: true,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("DefaultFormatWithVelocityLog", func(bb *testing.B) {
		app := velocity.New()
		logger := velocitylog.DefaultLogger()
		logger.SetOutput(io.Discard)
		app.Use(New(Config{
			Output: LoggerToWriter(logger, velocitylog.LevelDebug),
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("WithTagParameter", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status} ${reqHeader:test}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("test", "test")
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("WithLocals", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${locals:demo}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Locals("demo", "johndoe")
			return c.SendStatus(velocity.StatusOK)
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("WithLocalsInt", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${locals:demo}",
			Output: io.Discard,
		}))
		app.Get("/int", func(c velocity.Ctx) error {
			c.Locals("demo", 55)
			return c.SendStatus(velocity.StatusOK)
		})
		benchmarkSetup(bb, app, "/int")
	})

	b.Run("WithCustomDone", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Done: func(c velocity.Ctx, logString []byte) {
				if c.Response().StatusCode() == velocity.StatusOK {
					io.Discard.Write(logString) //nolint:errcheck // ignore error
				}
			},
			Output: io.Discard,
		}))
		app.Get("/logging", func(ctx velocity.Ctx) error {
			return ctx.SendStatus(velocity.StatusOK)
		})
		benchmarkSetup(bb, app, "/logging")
	})

	b.Run("WithAllTags", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${pid}${reqHeaders}${referer}${scheme}${protocol}${ip}${ips}${host}${url}${ua}${body}${route}${black}${red}${green}${yellow}${blue}${magenta}${cyan}${white}${reset}${error}${reqHeader:test}${query:test}${form:test}${cookie:test}${non}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("Streaming", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("Connection", "keep-alive")
			c.Set("Transfer-Encoding", "chunked")
			c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
				var i int
				for {
					i++
					msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
					fmt.Fprintf(w, "data: Message: %s\n\n", msg) //nolint:errcheck // ignore error
					err := w.Flush()
					if err != nil {
						break
					}
					if i == 10 {
						break
					}
				}
			})
			return nil
		})
		benchmarkSetup(bb, app, "/")
	})

	b.Run("WithBody", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${resBody}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Sample response body")
		})
		benchmarkSetup(bb, app, "/")
	})
}

// go test -v -run=^$ -bench=Benchmark_Logger_Parallel$ -benchmem -count=4
func Benchmark_Logger_Parallel(b *testing.B) {
	b.Run("NoMiddleware", func(bb *testing.B) {
		app := velocity.New()
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("WithBytesAndStatus", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("test", "test")
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("DefaultFormat", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("DefaultFormatWithVelocityLog", func(bb *testing.B) {
		app := velocity.New()
		logger := velocitylog.DefaultLogger()
		logger.SetOutput(io.Discard)
		app.Use(New(Config{
			Output: LoggerToWriter(logger, velocitylog.LevelDebug),
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("DefaultFormatDisableColors", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Output:        io.Discard,
			DisableColors: true,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("WithTagParameter", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status} ${reqHeader:test}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("test", "test")
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("WithLocals", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${locals:demo}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Locals("demo", "johndoe")
			return c.SendStatus(velocity.StatusOK)
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("WithLocalsInt", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${locals:demo}",
			Output: io.Discard,
		}))
		app.Get("/int", func(c velocity.Ctx) error {
			c.Locals("demo", 55)
			return c.SendStatus(velocity.StatusOK)
		})
		benchmarkSetupParallel(bb, app, "/int")
	})

	b.Run("WithCustomDone", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Done: func(c velocity.Ctx, logString []byte) {
				if c.Response().StatusCode() == velocity.StatusOK {
					io.Discard.Write(logString) //nolint:errcheck // ignore error
				}
			},
			Output: io.Discard,
		}))
		app.Get("/logging", func(ctx velocity.Ctx) error {
			return ctx.SendStatus(velocity.StatusOK)
		})
		benchmarkSetupParallel(bb, app, "/logging")
	})

	b.Run("WithAllTags", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${pid}${reqHeaders}${referer}${scheme}${protocol}${ip}${ips}${host}${url}${ua}${body}${route}${black}${red}${green}${yellow}${blue}${magenta}${cyan}${white}${reset}${error}${reqHeader:test}${query:test}${form:test}${cookie:test}${non}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Hello, World!")
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("Streaming", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${bytesReceived} ${bytesSent} ${status}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			c.Set("Connection", "keep-alive")
			c.Set("Transfer-Encoding", "chunked")
			c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
				var i int
				for {
					i++
					msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
					fmt.Fprintf(w, "data: Message: %s\n\n", msg) //nolint:errcheck // ignore error
					err := w.Flush()
					if err != nil {
						break
					}
					if i == 10 {
						break
					}
				}
			})
			return nil
		})
		benchmarkSetupParallel(bb, app, "/")
	})

	b.Run("WithBody", func(bb *testing.B) {
		app := velocity.New()
		app.Use(New(Config{
			Format: "${resBody}",
			Output: io.Discard,
		}))
		app.Get("/", func(c velocity.Ctx) error {
			return c.SendString("Sample response body")
		})
		benchmarkSetupParallel(bb, app, "/")
	})
}
