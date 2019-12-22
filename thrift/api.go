package thrift

type Status int64

type APIResponse struct {
	Status    Status            `thrift:"status,1" db:"status" json:"status"`
	Message   string            `thrift:"message,2" db:"message" json:"message"`
	Headers   map[string]string `thrift:"headers,3" db:"headers" json:"headers"`
	Content   string            `thrift:"content,4" db:"content" json:"content"`
	Total     int64             `thrift:"total,5" db:"total" json:"total"`
	ErrorCode string            `thrift:"errorCode,6" db:"errorCode" json:"errorCode"`
}
