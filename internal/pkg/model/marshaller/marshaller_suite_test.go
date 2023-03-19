package marshaller_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMarshaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Marshaller Suite")
}
