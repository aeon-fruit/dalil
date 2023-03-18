package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	"github.com/aeon-fruit/dalil.git/internal/pkg/middleware"
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
						value, err := middleware.GetPathParamInt(next.request.Context(), key)

						Expect(value).To(Equal(param))
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("PathParamContextString", func() {
		const customPattern = "[a-zA-Z]+-[0-9]+"
		var mw (func(http.Handler) http.Handler)

		for _, pattern := range []string{"", customPattern} {
			Context(fmt.Sprintf("the pattern is '%v'", pattern), func() {
				BeforeEach(func() {
					mw = middleware.PathParamContextString(key, pattern)
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
						if pattern != "" {
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
						}

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
								value, err := middleware.GetPathParamString(next.request.Context(), key)

								Expect(value).To(Equal(param))
								Expect(err).ToNot(HaveOccurred())
							})
						})
					})
				})
			})
		}
	})

	Describe("GetPathParamInt", func() {
		When("context argument is nil", func() {
			ctx := context.Context(nil)
			value, err := middleware.GetPathParamInt(ctx, key)

			It("returns 0 and error", func() {
				Expect(value).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})

		When("key does not match a value", func() {
			value, err := middleware.GetPathParamInt(context.TODO(), key)

			It("returns 0 and error", func() {
				Expect(value).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})
	})

	Describe("GetPathParamString", func() {
		When("context argument is nil", func() {
			ctx := context.Context(nil)
			value, err := middleware.GetPathParamString(ctx, key)

			It("returns 0 and error", func() {
				Expect(value).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})

		When("key does not match a value", func() {
			value, err := middleware.GetPathParamString(context.TODO(), key)

			It("returns 0 and error", func() {
				Expect(value).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})
	})
})
