package gin

import (
	"github.com/gin-gonic/gin"
)

type text3OutputHandler struct {}
func (*text3OutputHandler) Output(context *gin.Context) error {
	/* Your implementation */
	return nil
}

func text3Handler() OutputHandler {
	return &text3OutputHandler{}
}

func ExampleOutputHandler() {
	builder := NewMvcConfig().ToBuilder()
	handler := builder.WrapToGinHandler(text3Handler)
	// ginEngine.Get("/text-3", handler)

	_ = handler
}
