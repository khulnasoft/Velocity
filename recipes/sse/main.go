package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
)

var index = []byte(`<!DOCTYPE html>
<html>
<body>

<h1>SSE Messages</h1>
<div id="result"></div>

<script>
if(typeof(EventSource) !== "undefined") {
  var source = new EventSource("http://127.0.0.1:3000/sse");
  source.onmessage = function(event) {
    document.getElementById("result").innerHTML += event.data + "<br>";
  };
} else {
  document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
}
</script>

</body>
</html>
`)

func main() {
	// Velocity instance
	app := velocity.New()

	// CORS for external resources
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Cache-Control",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *velocity.Ctx) error {
		c.Response().Header.SetContentType(velocity.MIMETextHTMLCharsetUTF8)

		return c.Status(velocity.StatusOK).Send(index)
	})

	app.Get("/sse", func(c *velocity.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Status(velocity.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				fmt.Println(msg)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(2 * time.Second)
			}
		}))

		return nil
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
