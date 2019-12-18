package model

// MethodEnum ...
type MethodEnum struct {
	GET     *MethodValue
	POST    *MethodValue
	PUT     *MethodValue
	DELETE  *MethodValue
	OPTIONS *MethodValue
}

// APIMethod Published enum
var APIMethod = MethodEnum{
	GET:     &MethodValue{Value: "GET"},
	POST:    &MethodValue{Value: "POST"},
	PUT:     &MethodValue{Value: "PUT"},
	DELETE:  &MethodValue{Value: "DELETE"},
	OPTIONS: &MethodValue{Value: "OPTIONS"},
}

// StatusEnum ...
type StatusEnum struct {
	Ok           string
	Error        string
	Invalid      string
	NotFound     string
	Forbidden    string
	Existed      string
	Unauthorized string
}

// APIStatus Published enum
var APIStatus = &StatusEnum{
	Ok:           "OK",
	Error:        "ERROR",
	Invalid:      "INVALID",
	NotFound:     "NOT_FOUND",
	Forbidden:    "FORBIDDEN",
	Existed:      "EXISTED",
	Unauthorized: "UNAUTHORIZED",
}
