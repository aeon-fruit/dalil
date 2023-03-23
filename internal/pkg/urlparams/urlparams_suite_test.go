package urlparams_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestURLParams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "URLParams Suite")
}
