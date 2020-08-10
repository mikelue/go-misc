package gin

import (
	"github.com/gin-gonic/gin"
)

func text3OutputFunc(context *gin.Context) error {
	/* Your implementation */
	return nil
}

func text3HandlerByFunc() OutputHandler {
	return OutputHandlerFunc(text3OutputFunc)
}

func ExampleOutputHandlerFunc() {
	builder := NewMvcConfig().ToBuilder()
	handler := builder.WrapToGinHandler(text3Handler)
	// ginEngine.Get("/text-3", handler)

	_ = handler
}
