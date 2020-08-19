package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func text3OutputFunc(context *gin.Context) error {
	context.Data(http.StatusNotAcceptable, "plain/text", []byte("This is fine"))
	return nil
}

func text3HandlerByFunc() OutputHandler {
	return OutputHandlerFunc(text3OutputFunc)
}

func ExampleOutputHandlerFunc() {
	/**
	 * Prepares request
	 */
	sampleContext, resp := newContext()
	// :~)

	/**
	 * Wraps the customized handler
	 */
	handler := NewMvcConfig().ToBuilder().
		WrapToGinHandler(text3OutputFunc)
	handler(sampleContext)
	// :~)

	fmt.Printf("Resp[%d]: %s", resp.Code, resp.Body.String())
	// Output:
	// Resp[406]: This is fine
}
