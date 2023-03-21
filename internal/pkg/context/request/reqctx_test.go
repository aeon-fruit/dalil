package reqctx_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
)

var _ = Describe("ReqCtx", func() {

	const (
		key      = constants.Id
		intValue = 127
		value    = "127"
	)

	nilCtx := context.Context(nil)

	Describe("PathParam", func() {

		BeforeEach(func() {
			ctx := reqctx.SetPathParam(context.Background(), key, value)

			Expect(ctx).NotTo(BeZero())

			param, err := reqctx.GetPathParam(ctx, key)

			Expect(err).ToNot(HaveOccurred())
			Expect(param.String()).To(Equal(value))
			Expect(param.Int()).To(Equal(intValue))
		})

		Describe("Int", func() {

			When("an integer can be parsed from the underlying value", func() {
				It("returns the int value and no error", func() {})
			})

			When("an integer can be parsed from the underlying value", func() {
				It("returns 0 and an error ErrNotFound", func() {
					const nonInteger = "non-integer"
					ctx := reqctx.SetPathParam(context.Background(), key, nonInteger)

					Expect(ctx).NotTo(BeZero())

					param, err := reqctx.GetPathParam(ctx, key)

					Expect(err).ToNot(HaveOccurred())
					Expect(param.String()).To(Equal(nonInteger))

					intParam, err := param.Int()

					Expect(intParam).To(BeZero())
					Expect(err).To(Equal(errors.ErrNotFound))
				})
			})

		})

		Describe("String", func() {

			When("called", func() {
				It("returns the string value", func() {})
			})

		})

	})

	Describe("SetPathParam", func() {

		When("the context is nil", func() {
			It("returns a new context with the value enclosed in a PathParam", func() {
				ctx := reqctx.SetPathParam(nilCtx, key, value)

				Expect(ctx).NotTo(BeZero())

				param, err := reqctx.GetPathParam(ctx, key)

				Expect(err).ToNot(HaveOccurred())
				Expect(param.String()).To(Equal(value))
				Expect(param.Int()).To(Equal(intValue))
			})
		})

		When("the context is non nil", func() {

			var ctx context.Context

			BeforeEach(func() {
				ctx = reqctx.SetPathParam(context.Background(), key, value)

				Expect(ctx).NotTo(BeZero())

				param, err := reqctx.GetPathParam(ctx, key)

				Expect(err).ToNot(HaveOccurred())
				Expect(param.String()).To(Equal(value))
				Expect(param.Int()).To(Equal(intValue))
			})

			When("the key is not in the context", func() {
				It("returns a new context with the value enclosed in a PathParam", func() {})
			})

			When("the key is in the context", func() {
				It("returns the same context with the new value enclosed in a PathParam", func() {
					const (
						newIntVale = 10
						newValue   = "10"
					)

					newCtx := reqctx.SetPathParam(ctx, key, newValue)

					Expect(newCtx).To(Equal(ctx))

					param, err := reqctx.GetPathParam(ctx, key)

					Expect(err).ToNot(HaveOccurred())
					Expect(param.String()).To(Equal(newValue))
					Expect(param.Int()).To(Equal(newIntVale))
				})
			})

		})

	})

	Describe("GetPathParam", func() {

		When("context is nil", func() {
			It("returns an error ErrNotFound", func() {
				param, err := reqctx.GetPathParam(nilCtx, key)

				Expect(param).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})

		When("key not found", func() {
			It("returns an error ErrNotFound", func() {
				param, err := reqctx.GetPathParam(context.TODO(), "some not found key")

				Expect(param).To(BeZero())
				Expect(err).To(Equal(errors.ErrNotFound))
			})
		})

	})

})
