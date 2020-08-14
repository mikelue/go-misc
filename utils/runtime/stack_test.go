package runtime

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stack", func() {
	Context("CallerUtils", func() {
		It("GetDirOfSource", func() {
			testedDir := CallerUtils.GetDirOfSource()
			Expect(testedDir).To(ContainSubstring("go-misc/utils/runtime"))
		})
	})
})
