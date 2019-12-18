package model

import "github.com/SmartRice/gateway/thrift"

// APIResponder ...
type APIResponder interface {
	Respond(*APIResponse) error
	GetThriftResponse() *thrift.APIResponse
}

// APIResponse This is  response object with JSON format
type APIResponse struct {
	Status    string            `json:"status"`
	Data      interface{}       `json:"data,omitempty"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode,omitempty"`
	Total     int64             `json:"total,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}
