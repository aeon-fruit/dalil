package urlparams_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/urlparams"
)

var _ = Describe("URLParams", func() {

	const (
		key   = "id"
		value = "value"
	)

	Describe("ParsePathParam", func() {

		const value = "value"
		var request *http.Request

		BeforeEach(func() {
			request = httptest.NewRequest("", "http://url", strings.NewReader(""))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(key, value)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
		})

		When("the param is not found", func() {
			It("returns an empty string", func() {
				Expect(urlparams.ParsePathParam(request, key+"-suffix")).To(BeEmpty())
			})
		})

		When("the param is found", func() {
			It("returns its value", func() {
				Expect(urlparams.ParsePathParam(request, key)).To(Equal(value))
			})
		})

	})

	Describe("ParseQueryParam", func() {

		var request *http.Request

		BeforeEach(func() {
			request = httptest.NewRequest("", "http://url", strings.NewReader(""))
			query := request.URL.Query()
			query.Add(key, value)
			request.URL.RawQuery = query.Encode()
		})

		When("the param is not found", func() {
			It("returns an empty string", func() {
				Expect(urlparams.ParseQueryParam(request, key+"-suffix")).To(BeEmpty())
			})
		})

		When("the param is found", func() {
			It("returns its value", func() {
				Expect(urlparams.ParseQueryParam(request, key)).To(Equal(value))
			})
		})

	})

	Describe("ParseQueryFlag", func() {

		var (
			request *http.Request
			query   url.Values
		)

		BeforeEach(func() {
			request = httptest.NewRequest("", "http://url", strings.NewReader(""))
			query = request.URL.Query()
		})

		When("the param is not found", func() {
			It("returns false", func() {
				query.Add(key, "true")
				request.URL.RawQuery = query.Encode()

				Expect(urlparams.ParseQueryFlag(request, key+"-suffix")).To(BeFalse())
			})
		})

		When("the value of is convertible to false", func() {
			It("returns false", func() {
				query.Add(key, "0")
				request.URL.RawQuery = query.Encode()

				Expect(urlparams.ParseQueryFlag(request, key)).To(BeFalse())
			})
		})

		When("the value of is convertible to true", func() {
			It("returns true", func() {
				query.Add(key, "T")
				request.URL.RawQuery = query.Encode()

				Expect(urlparams.ParseQueryFlag(request, key)).To(BeTrue())
			})
		})

		When("the value of is empty", func() {
			It("returns true", func() {
				query.Add(key, "")
				request.URL.RawQuery = query.Encode()

				Expect(urlparams.ParseQueryFlag(request, key)).To(BeTrue())
			})
		})

		When("the value of is not empty and not convertible to a boolean", func() {
			It("returns false", func() {
				query.Add(key, value)
				request.URL.RawQuery = query.Encode()

				Expect(urlparams.ParseQueryFlag(request, key)).To(BeFalse())
			})
		})

	})

})
