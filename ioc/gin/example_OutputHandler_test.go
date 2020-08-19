package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type text3OutputHandler struct {}
func (*text3OutputHandler) Output(context *gin.Context) error {
	context.Data(http.StatusUseProxy, "plain/text", []byte("hello world"))
	return nil
}

func text3Handler() OutputHandler {
	return &text3OutputHandler{}
}

func ExampleOutputHandler() {
	/**
	 * Prepares request
	 */
	sampleContext, resp := newContext()
	// :~)

	/**
	 * Wraps the customized handler
	 */
	handler := NewMvcConfig().ToBuilder().
		WrapToGinHandler(text3Handler)
	handler(sampleContext)
	// :~)

	fmt.Printf("Resp[%d]: %s", resp.Code, resp.Body.String())
	// Output:
	// Resp[305]: hello world
}
