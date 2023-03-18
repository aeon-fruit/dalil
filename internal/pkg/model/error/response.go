package error

import (
	"net/http"
	"time"

	stubs "github.com/aeon-fruit/dalil.git/internal/pkg/stub/time"
)

type Response struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ResponseOption func(*Response)

func New(httpStatusCode int, options ...ResponseOption) Response {
	instance := Response{
		Code:      httpStatusCode,
		Message:   http.StatusText(httpStatusCode),
		Timestamp: time.Now(),
	}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return instance
}

func WithMessage(message string) ResponseOption {
	return func(response *Response) {
		if response != nil {
			response.Message = message
		}
	}
}

func UsingClock(clock stubs.Clock) ResponseOption {
	return func(response *Response) {
		if response != nil {
			response.Timestamp = clock.Now()
		}
	}
}
