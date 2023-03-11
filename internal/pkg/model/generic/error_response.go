package generic

import (
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorResponse(httpStatusCode int) ErrorResponse {
	return ErrorResponse{
		Code:      httpStatusCode,
		Message:   http.StatusText(httpStatusCode),
		Timestamp: time.Now(),
	}
}

func (er ErrorResponse) WithMessage(message string) ErrorResponse {
	er.Message = message
	er.Timestamp = time.Now()
	return er
}
