package time_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/stub/time"
)

var _ = Describe("Time", func() {

	When("Now called", func() {
		It("returns a non-zero time", func() {
			Expect(time.New().Now()).NotTo(BeZero())
		})
	})

})
