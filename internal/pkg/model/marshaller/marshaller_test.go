package marshaller_test

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	commonErrors "github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
	"github.com/aeon-fruit/dalil.git/internal/pkg/model/marshaller"
)

type testClock struct {
	time time.Time
}

func (tc testClock) Now() time.Time {
	return tc.time
}

type errorResponseRecorder struct {
	httptest.ResponseRecorder
}

func (rw *errorResponseRecorder) Write(_ []byte) (int, error) {
	return 0, errors.New("error")
}

func (rw *errorResponseRecorder) WriteString(_ string) (int, error) {
	return 0, errors.New("error")
}

var _ = Describe("Marshaller", func() {
	var recorder *httptest.ResponseRecorder

	Describe("SerializeEntity", func() {
		var recorder *httptest.ResponseRecorder
		entity := map[string]string{"key": "value"}
		serializedEntity := "{\"key\":\"value\"}\n"

		BeforeEach(func() {
			recorder = httptest.NewRecorder()
		})

		When("writer argument is nil", func() {
			It("returns an error", func() {
				err := marshaller.SerializeEntity(nil, entity)

				Expect(err).To(HaveOccurred())
				Expect(errors.Unwrap(err)).To(Equal(commonErrors.ErrInvalidArgument))
			})
		})

		When("operation fails", func() {
			It("doesn't write the entity to the writer and returns an error", func() {
				err := marshaller.SerializeEntity(recorder, []any{math.Inf(-1)})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("json"))

				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Body.Bytes()).To(BeNil())
			})
		})

		When("operation succeeds", func() {
			It("writes the entity to the writer and doesn't return an error", func() {
				err := marshaller.SerializeEntity(recorder, entity)

				Expect(err).ToNot(HaveOccurred())

				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(Equal(serializedEntity))
			})
		})
	})

	Describe("SerializeError", func() {
		const statusCode = http.StatusTeapot
		timestamp := time.UnixMilli(1679143523911)
		clock := testClock{
			time: timestamp,
		}
		response := errorModel.New(statusCode, "", errorModel.UsingClock(clock))
		serializedEntity := "{\"code\":418,\"message\":\"I'm a teapot\",\"timestamp\":\"2023-03-18T13:45:23.911+01:00\"}\n"

		BeforeEach(func() {
			recorder = httptest.NewRecorder()
		})

		When("writer argument is nil", func() {
			It("returns an error", func() {
				err := marshaller.SerializeError(nil, response)

				Expect(err).To(HaveOccurred())
				Expect(errors.Unwrap(err)).To(Equal(commonErrors.ErrInvalidArgument))
			})
		})

		When("errorResponse argument is empty", func() {
			It("writes the errorResponse and StatusInternalServerError to the writer and doesn't return an error", func() {
				err := marshaller.SerializeError(recorder, errorModel.Response{})

				Expect(err).ToNot(HaveOccurred())

				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(Equal("{\"code\":0,\"message\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}\n"))
			})
		})

		When("operation fails", func() {
			It("doesn't write the errorResponse to the writer, set http status to StatusInternalServerError and returns an error", func() {
				err := marshaller.SerializeError(&errorResponseRecorder{*recorder}, response)

				Expect(err).To(HaveOccurred())

				Expect(recorder.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				Expect(recorder.Body.Bytes()).To(BeNil())
			})
		})

		When("operation succeeds", func() {
			It("writes the errorResponse to the writer and doesn't return an error", func() {
				err := marshaller.SerializeError(recorder, response)

				Expect(err).ToNot(HaveOccurred())

				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Result().StatusCode).To(Equal(statusCode))
				Expect(recorder.Body.String()).To(Equal(serializedEntity))
			})
		})
	})
})
