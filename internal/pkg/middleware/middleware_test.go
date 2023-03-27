package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
	"github.com/aeon-fruit/dalil.git/internal/pkg/middleware"
	logMock "github.com/aeon-fruit/dalil.git/test/mocks/log"
)

type testHandler struct {
	request   *http.Request
	callCount int
}

func (th *testHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	th.request = r
	th.callCount++
}

var _ = Describe("Middleware", func() {

	const key = "id"

	Describe("LoggingContext", func() {
		var mockLogger *logMock.MockLogger
		var mw (func(http.Handler) http.Handler)

		BeforeEach(func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockLogger = logMock.NewMockLogger(mockCtrl)
			mw = middleware.LoggingContext(mockLogger)
		})

		It("returns a non-nil function", func() {
			Expect(mw).NotTo(BeNil())
		})

		It("has a return that returns a non-nil http.Handler", func() {
			next := &testHandler{}
			handler := mw(next)

			Expect(handler).NotTo(BeNil())
		})

		Describe("http.Handler returned by the returned middleware", func() {
			var next *testHandler
			var handler http.Handler

			BeforeEach(func() {
				next = &testHandler{}
				handler = mw(next)

				Expect(handler).NotTo(BeNil())
			})

			When("used", func() {
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest("", "http://url", strings.NewReader(""))

				BeforeEach(func() {
					mockLogger.EXPECT().WithName(constants.AppName).Times(1)
					handler.ServeHTTP(recorder, request)
				})

				It("augments the request context by a logger named after the app", func() {
					Expect(next.callCount).NotTo(BeZero())
					Expect(next.request).NotTo(BeNil())

					Expect(logr.FromContext(next.request.Context())).NotTo(BeNil())
				})
			})

		})

	})

	Describe("PathParamContextInt", func() {
		var mw (func(http.Handler) http.Handler)

		BeforeEach(func() {
			mw = middleware.PathParamContextInt(key)
		})

		It("returns a non-nil function", func() {
			Expect(mw).NotTo(BeNil())
		})

		It("has a return that returns a non-nil http.Handler", func() {
			next := &testHandler{}
			handler := mw(next)

			Expect(handler).NotTo(BeNil())
		})

		Describe("http.Handler returned by the returned middleware", func() {
			var next *testHandler
			var handler http.Handler

			BeforeEach(func() {
				next = &testHandler{}
				handler = mw(next)

				Expect(handler).NotTo(BeNil())
			})

			When("key has no matching key", func() {
				recorder := httptest.NewRecorder()

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add(key+"-suffix", "non integer value")

				request := httptest.NewRequest("", "http://url", strings.NewReader(""))
				request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

				BeforeEach(func() {
					handler.ServeHTTP(recorder, request)
				})

				It("sends StatusBadRequest", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("doesn't forward the request to the next handler", func() {
					Expect(next.callCount).To(BeZero())
				})
			})

			Context("key has a matching value", func() {
				When("an integer cannot be parsed from the value", func() {
					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key, "non integer value")

					request := httptest.NewRequest("", "http://url", strings.NewReader(""))
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("sends StatusBadRequest", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})

					It("doesn't forward the request to the next handler", func() {
						Expect(next.callCount).To(BeZero())
					})
				})
				When("an integer can be parsed from the value", func() {
					const param = 12345

					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key, strconv.Itoa(param))

					request := (&http.Request{})
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("forwards the request to the next handler", func() {
						Expect(next.callCount).NotTo(BeZero())
					})

					It("augments the context with the parsed int from the value", func() {
						value, err := reqctx.GetPathParam(next.request.Context(), key)

						Expect(err).ToNot(HaveOccurred())

						intValue, err := value.Int()
						Expect(err).ToNot(HaveOccurred())
						Expect(intValue).To(Equal(param))
					})
				})
			})
		})
	})

	Describe("PathParamContextString", func() {
		const (
			customPattern = "[a-zA-Z]+-[0-9]+"
			errorPattern  = "[a-zA-Z]+-[0-9+"
		)

		Context("the pattern is empty", func() {
			var mw (func(http.Handler) http.Handler)

			BeforeEach(func() {
				mw = middleware.PathParamContextString(key, "")
			})

			It("returns a non-nil function", func() {
				Expect(mw).NotTo(BeNil())
			})

			It("has a return that returns a non-nil http.Handler", func() {
				next := &testHandler{}
				handler := mw(next)

				Expect(handler).NotTo(BeNil())
			})

			Describe("http.Handler returned by the returned middleware", func() {
				var next *testHandler
				var handler http.Handler

				BeforeEach(func() {
					next = &testHandler{}
					handler = mw(next)

					Expect(handler).NotTo(BeNil())
				})

				When("key has no matching key", func() {
					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key+"-suffix", "non integer value")

					request := httptest.NewRequest("", "http://url", strings.NewReader(""))
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("sends StatusBadRequest", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})

					It("doesn't forward the request to the next handler", func() {
						Expect(next.callCount).To(BeZero())
					})
				})

				Context("key has a matching value", func() {
					When("the value matches the pattern", func() {
						const param = "Match-12345"

						recorder := httptest.NewRecorder()

						rctx := chi.NewRouteContext()
						rctx.URLParams.Add(key, param)

						request := (&http.Request{})
						request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

						BeforeEach(func() {
							handler.ServeHTTP(recorder, request)
						})

						It("forwards the request to the next handler", func() {
							Expect(next.callCount).NotTo(BeZero())
						})

						It("augments the context with the value", func() {
							value, err := reqctx.GetPathParam(next.request.Context(), key)

							Expect(value.String()).To(Equal(param))
							Expect(err).ToNot(HaveOccurred())
						})
					})
				})
			})
		})

		Context("the pattern is a non-empty valid regex", func() {
			var mw (func(http.Handler) http.Handler)

			BeforeEach(func() {
				mw = middleware.PathParamContextString(key, customPattern)
			})

			It("returns a non-nil function", func() {
				Expect(mw).NotTo(BeNil())
			})

			It("has a return that returns a non-nil http.Handler", func() {
				next := &testHandler{}
				handler := mw(next)

				Expect(handler).NotTo(BeNil())
			})

			Describe("http.Handler returned by the returned middleware", func() {
				var next *testHandler
				var handler http.Handler

				BeforeEach(func() {
					next = &testHandler{}
					handler = mw(next)

					Expect(handler).NotTo(BeNil())
				})

				When("key has no matching key", func() {
					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key+"-suffix", "non integer value")

					request := httptest.NewRequest("", "http://url", strings.NewReader(""))
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("sends StatusBadRequest", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})

					It("doesn't forward the request to the next handler", func() {
						Expect(next.callCount).To(BeZero())
					})
				})

				Context("key has a matching value", func() {
					When("the value doesn't match the pattern", func() {
						recorder := httptest.NewRecorder()

						rctx := chi.NewRouteContext()
						rctx.URLParams.Add(key, "non-matching")

						request := httptest.NewRequest("", "http://url", strings.NewReader(""))
						request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

						BeforeEach(func() {
							handler.ServeHTTP(recorder, request)
						})

						It("sends StatusBadRequest", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})

						It("doesn't forward the request to the next handler", func() {
							Expect(next.callCount).To(BeZero())
						})
					})

					When("the value matches the pattern", func() {
						const param = "Match-12345"

						recorder := httptest.NewRecorder()

						rctx := chi.NewRouteContext()
						rctx.URLParams.Add(key, param)

						request := (&http.Request{})
						request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

						BeforeEach(func() {
							handler.ServeHTTP(recorder, request)
						})

						It("forwards the request to the next handler", func() {
							Expect(next.callCount).NotTo(BeZero())
						})

						It("augments the context with the value", func() {
							value, err := reqctx.GetPathParam(next.request.Context(), key)

							Expect(value.String()).To(Equal(param))
							Expect(err).ToNot(HaveOccurred())
						})
					})
				})
			})
		})

		Context("the pattern is not a valid regex", func() {
			var mw (func(http.Handler) http.Handler)

			BeforeEach(func() {
				mw = middleware.PathParamContextString(key, errorPattern)
			})

			It("returns a non-nil function", func() {
				Expect(mw).NotTo(BeNil())
			})

			It("has a return that returns a non-nil http.Handler", func() {
				next := &testHandler{}
				handler := mw(next)

				Expect(handler).NotTo(BeNil())
			})

			Describe("http.Handler returned by the returned middleware", func() {
				var next *testHandler
				var handler http.Handler

				BeforeEach(func() {
					next = &testHandler{}
					handler = mw(next)

					Expect(handler).NotTo(BeNil())
				})

				When("key has no matching key", func() {
					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key+"-suffix", "non integer value")

					request := httptest.NewRequest("", "http://url", strings.NewReader(""))
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("sends StatusBadRequest", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})

					It("doesn't forward the request to the next handler", func() {
						Expect(next.callCount).To(BeZero())
					})
				})

				Context("key has a matching value", func() {
					recorder := httptest.NewRecorder()

					rctx := chi.NewRouteContext()
					rctx.URLParams.Add(key, "non-matching")

					request := httptest.NewRequest("", "http://url", strings.NewReader(""))
					request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

					BeforeEach(func() {
						handler.ServeHTTP(recorder, request)
					})

					It("sends StatusBadRequest", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})

					It("doesn't forward the request to the next handler", func() {
						Expect(next.callCount).To(BeZero())
					})
				})
			})
		})

	})

})
