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
		It("has a code and a message equal to the arguments, and a non zero timestamp", func() {
			const (
				code    = 9000
				message = "an error message"
			)
			instance := error.New(code, message)

			Expect(instance.Code).To(Equal(code))
			Expect(instance.Message).To(Equal(message))
			Expect(instance.Timestamp).NotTo(BeZero())
		})

		Context("message argument is empty", func() {
			When("code argument is a valid HTTP status code", func() {
				It("has a message initialized with the text of the HTTP status", func() {
					const httpStatus = http.StatusOK
					instance := error.New(httpStatus, "")

					Expect(instance.Message).To(Equal(http.StatusText(httpStatus)))
				})
			})

			When("code argument is not a valid HTTP status code", func() {
				It("has an empty message", func() {
					const invalidHttpStatus = 9000
					instance := error.New(invalidHttpStatus, "")

					Expect(instance.Message).To(BeEmpty())
				})
			})
		})

		Context("message argument is not empty", func() {
			const message = "an error message"
			When("code argument is a valid HTTP status code", func() {
				It("has a message initialized with the message argument", func() {
					const httpStatus = http.StatusOK
					instance := error.New(httpStatus, message)

					Expect(instance.Message).To(Equal(message))
				})
			})

			When("code argument is not a valid HTTP status code", func() {
				It("has a message initialized with the message argument", func() {
					const invalidHttpStatus = 9000
					instance := error.New(invalidHttpStatus, message)

					Expect(instance.Message).To(Equal(message))
				})
			})
		})

		When("WithClock is specified", func() {
			It("has a current timestamp", func() {
				timestamp := time.UnixMilli(1679143523911)
				clock := testClock{
					time: timestamp,
				}

				instance := error.New(http.StatusConflict, "", error.UsingClock(clock))

				Expect(instance.Timestamp).To(Equal(timestamp))
			})
		})
	})
})
