package slf4go_logrus

import (
	"io"
	"strings"
	l4 "github.com/go-eden/slf4go"
	lr "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("logrusDriver", func() {
	Context("Logged Content with level setting", func() {
		var stringCapture *strings.Builder

		BeforeEach(func() {
			stringCapture = new(strings.Builder)

			UseLogrus.WithConfig(LogrousConfig{
				DEFAULT_LOGGER: newLogrusLoggerForTest(lr.InfoLevel, stringCapture),
				"log.s1": newLogrusLoggerForTest(lr.DebugLevel, stringCapture),
			})
		})

		It("Uses default logger(debug is skipped)", func() {
			l4.NewLogger("log.g1").Info("hello-1")
			l4.NewLogger("log.g1").Debug("hello-2")

			logResult := stringCapture.String()

			Expect(logResult).To(ContainSubstring("hello-1"))
			Expect(logResult).ToNot(ContainSubstring("hello-2"))
		})

		It("Uses named writer(debug is logged out)", func() {
			l4.NewLogger("log.g1").Info("hello-1")
			l4.NewLogger("log.g1").Debug("hello-2")
			l4.NewLogger("log.s1").Debug("hello-3")
			l4.NewLogger("log.s1").Trace("hello-4")

			logResult := stringCapture.String()

			Expect(logResult).To(And(ContainSubstring("hello-1"), ContainSubstring("hello-3")))
			Expect(logResult).ToNot(Or(ContainSubstring("hello-2"), ContainSubstring("hello-4")))
		})
	})

	Context("Internal keeper of logger", func() {
		var testedDriver *logrusDriver

		BeforeEach(func() {
			testedDriver = newLogrusDriver()
			testedDriver.loggerMap["test.log.r1"] = lr.New()
			testedDriver.loggerMap["test.log.r1"].SetLevel(lr.InfoLevel)
			testedDriver.loggerLevelMap["test.log.r1"] = l4.InfoLevel
		})

		DescribeTable("getLogger",
			func(name string, expectedOk bool) {
				testedLogger, ok := testedDriver.getLogger(name)

				Expect(ok).To(BeEquivalentTo(expectedOk))
				if expectedOk {
					Expect(testedLogger.GetLevel()).To(BeEquivalentTo(lr.InfoLevel))
				}
			},
			Entry("Existing", "test.log.r1", true),
			Entry("Not existing", "test.log.r2", false),
		)

		DescribeTable("getLevel",
			func(name string, expectedOk bool) {
				testedLevel, ok := testedDriver.getLevel(name)

				Expect(ok).To(BeEquivalentTo(expectedOk))
				if expectedOk {
					Expect(testedLevel).To(BeEquivalentTo(l4.InfoLevel))
				}
			},
			Entry("Existing", "test.log.r1", true),
			Entry("Not existing", "test.log.r2", false),
		)
	})
})

func newLogrusLoggerForTest(level lr.Level, writer io.Writer) *lr.Logger {
	newLogger := lr.New()
	newLogger.SetLevel(level)
	newLogger.SetOutput(writer)
	return newLogger
}
