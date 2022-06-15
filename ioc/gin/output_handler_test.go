package gin

import (
	"net/http"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("OutputHandler", func() {
	Context("builderByAccept", func() {
		It("Nil handler", func() {
			context := newContextByMime("image/nothing")

			Expect(builderByAccept(context)).To(BeNil())
		})

		DescribeTable("Viable handler",
			func(sampleAccepted string, expectedBuilder outputHandlerBuilder) {
				context := newContextByMime(sampleAccepted)

				testedAddress := fmt.Sprintf("%p", builderByAccept(context))
				expectedAddress := fmt.Sprintf("%p", expectedBuilder)
				Expect(testedAddress).To(BeEquivalentTo(expectedAddress))
			},
			Entry("json", "application/json", JsonOutputHandler),
			Entry("xml", "application/xml", XmlOutputHandler),
			Entry("xml(text)", "text/xml", XmlOutputHandler),
			Entry("text", "text/plain", TextOutputHandler),
			Entry("yaml", "application/x-yaml", YamlOutputHandler),
			Entry("protobuf", "application/x-protobuf", ProtoBufOutputHandler),
		)
	})

	Context("Output Handlers", itForHandlers)
})

var itForHandlers = func() {
	DescribeTable("XXXOutputHandler",
		func(
			sampleStatus int, sampleBody string,
			testedBuilder outputHandlerBuilder,
			expectedContentType string,
		) {
			context, testedRespRecorder := newContext()
			testedBuilder(sampleStatus, sampleBody).Output(context)

			testedResp := testedRespRecorder.Result()
			Expect(testedResp.StatusCode).To(BeEquivalentTo(sampleStatus))
			Expect(testedResp.Header.Get("Content-Type")).To(ContainSubstring(expectedContentType))
			Eventually(BufferReader(testedResp.Body)).
				Should(Say(sampleBody))
		},
		Entry("JsonOutputHandler", http.StatusFound, "[11, 21]", JsonOutputHandler, "application/json"),
		Entry("TextOutputHandler", http.StatusOK, "Hello World!", TextOutputHandler, "text/plain"),
		Entry("XmlOutputHandler", http.StatusConflict, "Hello", XmlOutputHandler, "application/xml"),
		Entry("YamlOutputHandler", http.StatusCreated, "{ a: 20, b: 40 }", YamlOutputHandler, "application/x-yaml"),
	)

	It("ProtoBufOutputHandler", func() {
		context, testedRespRecorder := newContext()
		sampleProtobuf := newPanda("Burton", 23)
		ProtoBufOutputHandler(http.StatusOK, sampleProtobuf).Output(context)

		testedResp := testedRespRecorder.Result()
		Expect(testedResp.StatusCode).To(BeEquivalentTo(http.StatusOK))
		Expect(testedResp.Header.Get("Content-Type")).To(ContainSubstring("application/x-protobuf"))
		Eventually(BufferReader(testedResp.Body)).Should(Say("Burton"))
	})
}

type panda struct {
	Label *string `protobuf:"bytes,1,req,name=label"`
	Type *int32 `protobuf:"varint,2,opt,name=type,def=77"`
}
func newPanda(label string, pandaType int32) *panda {
	newPanda := &panda{}
	newPanda.Label = &label
	newPanda.Type = &pandaType
	return newPanda
}
func (self *panda) Reset() { *self = panda{} }
func (self *panda) String() string { return fmt.Sprintf("panda[%s][%d]", *self.Label, *self.Type) }
func (*panda) ProtoMessage() {}
