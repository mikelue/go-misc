package runtime

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runtime Suite")
}
