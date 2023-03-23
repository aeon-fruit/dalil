package error

import (
	"net/http"
	"strings"
	"time"

	stubs "github.com/aeon-fruit/dalil.git/internal/pkg/stub/time"
)

type Response struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ResponseOption func(*Response)

func New(httpStatusCode int, message string, options ...ResponseOption) Response {
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(httpStatusCode)
	}

	instance := Response{
		Code:      httpStatusCode,
		Message:   message,
		Timestamp: time.Now(),
	}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return instance
}

func UsingClock(clock stubs.Clock) ResponseOption {
	return func(response *Response) {
		if response != nil {
			response.Timestamp = clock.Now()
		}
	}
}
