package utils

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mrt "github.com/mikelue/go-misc/utils/runtime"
)

var testSourceDir = mrt.CallerUtils.GetDirOfSource()

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}
