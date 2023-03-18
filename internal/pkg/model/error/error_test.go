package error_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
)

type testClock struct {
	time time.Time
}

func (tc testClock) Now() time.Time {
	return tc.time
}

var _ = Describe("Error", func() {
	Describe("New", func() {
		It("has a code equal to the argument and a non zero timestamp", func() {
			const value = 9000
			instance := error.New(value)

			Expect(instance.Code).To(Equal(value))
			Expect(instance.Timestamp).NotTo(BeZero())
		})

		When("argument is a valid HTTP status code", func() {
			It("has a message initialized with the text of the HTTP status", func() {
				const httpStatus = http.StatusOK
				instance := error.New(httpStatus)

				Expect(instance.Message).To(Equal(http.StatusText(httpStatus)))
			})
		})

		When("argument is not a valid HTTP status code", func() {
			It("has an empty message", func() {
				const invalidHttpStatus = 9000
				instance := error.New(invalidHttpStatus)

				Expect(instance.Message).To(BeEmpty())
			})
		})

		When("WithMessage is specified", func() {
			It("has a message initialized with the specified text", func() {
				const text = "The new message body"
				instance := error.New(http.StatusOK, error.WithMessage(text))

				Expect(instance.Message).To(Equal(text))
			})
		})

		When("WithClock is specified", func() {
			It("has a current timestamp", func() {
				timestamp := time.UnixMilli(1679143523911)
				clock := testClock{
					time: timestamp,
				}

				const httpStatus = http.StatusConflict
				instance := error.New(httpStatus, error.UsingClock(clock))

				Expect(instance.Timestamp).To(Equal(timestamp))
			})
		})
	})
})
