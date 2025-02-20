// Package adapter ...
// author : @fenny (author of velocity)
// Package velocityadapter adds Velocity support for the aws-severless-go-api library.
// Uses the core package behind the scenes and exposes the New method to
// get a new instance and Proxy method to send request to the Velocity app.
package adapter

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/khulnasoft/utils"
	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
)

// VelocityLambda makes it easy to send API Gateway proxy events to a velocity.App.
// The library transforms the proxy event into an HTTP request and then
// creates a proxy response object from the *velocity.Ctx
type VelocityLambda struct {
	core.RequestAccessor
	app *velocity.App
}

// New creates a new instance of the VelocityLambda object.
// Receives an initialized *velocity.App object - normally created with velocity.New().
// It returns the initialized instance of the VelocityLambda object.
func New(app *velocity.App) *VelocityLambda {
	return &VelocityLambda{
		app: app,
	}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the velocity.App for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (f *VelocityLambda) Proxy(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	velocityRequest, err := f.ProxyEventToHTTPRequest(req)
	return f.proxyInternal(velocityRequest, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the echo.Echo for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (f *VelocityLambda) ProxyWithContext(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	velocityRequest, err := f.EventToRequestWithContext(ctx, req)
	return f.proxyInternal(velocityRequest, err)
}

func (f *VelocityLambda) proxyInternal(req *http.Request, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	resp := core.NewProxyResponseWriter()
	f.adaptor(resp, req)

	proxyResponse, err := resp.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}

func (f *VelocityLambda) adaptor(w http.ResponseWriter, r *http.Request) {
	// New fasthttp request
	var req fasthttp.Request
	// Convert net/http -> fasthttp request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, utils.StatusMessage(velocity.StatusInternalServerError), velocity.StatusInternalServerError)
		return
	}
	req.Header.SetMethod(r.Method)
	req.SetRequestURI(r.RequestURI)
	req.Header.SetContentLength(len(body))
	req.SetHost(r.Host)
	for key, val := range r.Header {
		for _, v := range val {
			req.Header.Add(key, v)
		}
	}
	_, _ = req.BodyWriter().Write(body)
	remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		http.Error(w, utils.StatusMessage(velocity.StatusInternalServerError), velocity.StatusInternalServerError)
		return
	}

	// New fasthttp Ctx
	var fctx fasthttp.RequestCtx
	fctx.Init(&req, remoteAddr, nil)

	// Pass RequestCtx to Velocity router
	f.app.Handler()(&fctx)
	// Convert fasthttp Ctx > net/http
	fctx.Response.Header.VisitAll(func(k, v []byte) {
		sk := string(k)
		sv := string(v)
		w.Header().Set(sk, sv)
	})
	w.WriteHeader(fctx.Response.StatusCode())
	_, _ = w.Write(fctx.Response.Body())
}
