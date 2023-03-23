package repository_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	repository "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao"
)

var _ = Describe("Repository", func() {

	When("used", func() {
		It("needs Ginkgo specs", func() {
			repo := repository.New()
			Expect(repo).NotTo(BeNil())
		})
	})

})
