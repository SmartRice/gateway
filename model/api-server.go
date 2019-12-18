package model

import (
	"errors"
	"fmt"
	"github.com/SmartRice/gateway/thrift"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type APIServer interface {
	PreRequest(Handler) error
	SetHandler(*MethodValue, string, Handler) error
	Expose(int)
	Start(*sync.WaitGroup)
	GetHostname() string
}

// MethodValue ...
type MethodValue struct {
	Value string
}

type Handler = func(req APIRequest, res APIResponder) error

// HTTPAPIServer ...
type HTTPAPIServer struct {
	T        string
	Echo     *echo.Echo
	Port     int
	ID       int
	RunSSL   bool
	SSLPort  int
	hostname string
}

// HandlerWrapper handler object
type HandlerWrapper struct {
	handler Handler
	server  *HTTPAPIServer
}

// PreHandlerWrapper
type PreHandlerWrapper struct {
	preHandler Handler
	next       echo.HandlerFunc
	server     *HTTPAPIServer
}

// HTTPAPIResponder This is response object with JSON format
type HTTPAPIResponder struct {
	t        string
	context  echo.Context
	start    time.Time
	hostname string
}

func newHTTPAPIResponder(c echo.Context, hostname string) APIResponder {
	return &HTTPAPIResponder{
		t:        "HTTP",
		start:    time.Now(),
		context:  c,
		hostname: hostname,
	}
}

func (server *HTTPAPIServer) PreRequest(fn Handler) error {
	var preWrapper = &PreHandlerWrapper{
		preHandler: fn,
		server:     server,
	}

	server.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		preWrapper.next = next
		return preWrapper.processCore
	})
	return nil
}

func (server *HTTPAPIServer) SetHandler(method *MethodValue, path string, fn Handler) error {
	var wrapper = &HandlerWrapper{
		handler: fn,
		server:  server,
	}

	switch method.Value {
	case APIMethod.GET.Value:
		server.Echo.GET(path, wrapper.processCore)
	case APIMethod.POST.Value:
		server.Echo.POST(path, wrapper.processCore)
	case APIMethod.PUT.Value:
		server.Echo.PUT(path, wrapper.processCore)
	case APIMethod.DELETE.Value:
		server.Echo.DELETE(path, wrapper.processCore)
	}

	return nil
}

func (server *HTTPAPIServer) Expose(port int) {
	server.Port = port
}

func (server *HTTPAPIServer) Start(wg *sync.WaitGroup) {
	var ps = strconv.Itoa(server.Port)
	fmt.Println("  [ API Server " + strconv.Itoa(server.ID) + " ] Try to listen at " + ps)
	server.Echo.HideBanner = true

	if server.RunSSL {
		go server.Echo.StartTLS(":"+strconv.Itoa(server.SSLPort), "crt.pem", "key.pem")
	}

	server.Echo.Start(":" + ps)
	wg.Done()
}

func (server *HTTPAPIServer) GetHostname() string {
	return server.hostname
}

func NewHTTPAPIServer(id int, hostname string) APIServer {
	var server = HTTPAPIServer{
		T:        "HTTP",
		Echo:     echo.New(),
		ID:       id,
		hostname: hostname,
	}
	server.Echo.Use(middleware.Gzip())
	return &server
}

func (resp HTTPAPIResponder) Respond(response *APIResponse) error {
	var context = resp.context

	if response.Data != nil && reflect.TypeOf(response.Data).Kind() != reflect.Slice {
		return errors.New("data response must be a slice")
	}

	if response.Headers != nil {
		header := context.Response().Header()
		for key, value := range response.Headers {
			header.Set(key, value)
		}
		response.Headers = nil
	}

	var dif = float64(time.Since(resp.start).Nanoseconds()) / 1000000
	context.Response().Header().Set("X-Execution-Time", fmt.Sprintf("%.4f ms", dif))
	context.Response().Header().Set("X-Hostname", resp.hostname)

	switch response.Status {
	case APIStatus.Ok:
		return context.JSON(http.StatusOK, response)
	case APIStatus.Error:
		return context.JSON(http.StatusInternalServerError, response)
	case APIStatus.Forbidden:
		return context.JSON(http.StatusForbidden, response)
	case APIStatus.Invalid:
		return context.JSON(http.StatusBadRequest, response)
	case APIStatus.NotFound:
		return context.JSON(http.StatusNotFound, response)
	case APIStatus.Unauthorized:
		return context.JSON(http.StatusUnauthorized, response)
	case APIStatus.Existed:
		return context.JSON(http.StatusConflict, response)
	}

	return context.JSON(http.StatusBadRequest, response)
}

func (H HTTPAPIResponder) GetThriftResponse() *thrift.APIResponse {
	return nil
}

// processCore Process basic logic of Echo
func (hw *PreHandlerWrapper) processCore(c echo.Context) error {
	req := newHTTPAPIRequest(c)
	resp := newHTTPAPIResponder(c, hw.server.GetHostname())
	err := hw.preHandler(req, resp)
	if err == nil {
		hw.next(c)
	}
	return nil
}

// processCore Process basic logic of Echo
func (hw *HandlerWrapper) processCore(c echo.Context) error {
	hw.handler(newHTTPAPIRequest(c), newHTTPAPIResponder(c, hw.server.GetHostname()))
	return nil
}
