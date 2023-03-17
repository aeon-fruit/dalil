package error

import (
	"net/http"
	"time"
)

type Response struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func New(httpStatusCode int) Response {
	return Response{
		Code:      httpStatusCode,
		Message:   http.StatusText(httpStatusCode),
		Timestamp: time.Now(),
	}
}

func (er Response) WithMessage(message string) Response {
	er.Message = message
	er.Timestamp = time.Now()
	return er
}
