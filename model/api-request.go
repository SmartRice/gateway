package model

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
)

// APIRequest ...
type APIRequest interface {
	GetPath() string
	GetMethod() *MethodValue
	GetParam(string) string
	GetParams() map[string]string
	GetHeader(string) string
	GetHeaders() map[string]string
	GetContent(interface{}) error
	GetContentText() string
	GetAttribute(string) interface{}
	SetAttribute(string, interface{})
}

// HTTPAPIRequest This is response object with JSON format
type HTTPAPIRequest struct {
	t       string
	context echo.Context
	body    string
}

func (req *HTTPAPIRequest) GetParam(name string) string {
	return req.context.QueryParam(name)
}

func (req *HTTPAPIRequest) GetParams() map[string]string {
	var vals = req.context.QueryParams()
	var m = make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

func (req *HTTPAPIRequest) GetHeader(name string) string {
	return req.context.Request().Header.Get(name)
}

func (req *HTTPAPIRequest) GetHeaders() map[string]string {
	var vals = req.context.Request().Header
	var m = make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

func (req *HTTPAPIRequest) GetContent(data interface{}) error {
	return json.Unmarshal([]byte(req.GetContentText()), data)
}

func (req *HTTPAPIRequest) GetContentText() string {
	if req.body == "" {
		var bodyBytes []byte
		if req.context.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(req.context.Request().Body)
		}

		req.body = string(bodyBytes)
	}

	return req.body
}

func (req *HTTPAPIRequest) GetAttribute(name string) interface{} {
	return req.context.Get(name)
}

func (req *HTTPAPIRequest) SetAttribute(name string, value interface{}) {
	req.context.Set(name, value)
}

func newHTTPAPIRequest(e echo.Context) APIRequest {
	return &HTTPAPIRequest{
		t:       "HTTP",
		context: e,
	}
}

func (req *HTTPAPIRequest) GetPath() string {
	return req.context.Path()
}

func (req *HTTPAPIRequest) GetMethod() *MethodValue {
	var s = req.context.Request().Method
	switch s {
	case "GET":
		return APIMethod.GET
	case "POST":
		return APIMethod.POST
	case "PUT":
		return APIMethod.PUT
	case "DELETE":
		return APIMethod.DELETE
	}

	return &MethodValue{Value: s}
}
