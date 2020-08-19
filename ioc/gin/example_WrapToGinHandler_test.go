package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleMvcBuilder_wrapToGinHandler() {
	/**
	 * Prepares request
	 */
	req := httptest.NewRequest(
		"POST", "/sample-uri",
		strings.NewReader(`{ "id": 281, "color": 3 , "code": "IU-762"}`),
	)

	req.Header.Set("Content-Type", "application/json")

	sampleContext, resp := newContext()
	sampleContext.Request = req
	// :~)

	/**
	 * Wraps the customized handler
	 */
	handler := NewMvcConfig().ToBuilder().
		WrapToGinHandler(sampleHandler)
	handler(sampleContext)
	// :~)

	fmt.Printf("Resp[%d]: %s", resp.Code, resp.Body.String())
	// Output:
	// Resp[200]: [10,20,"IU-762"]
}

type leopard struct {
	Id int `json:"id"`
	Color uint8 `json:"color"`
	Code string `json:"code"`
}

func sampleHandler(leopard *leopard) OutputHandler {
	sampleOutput := []interface{} {
		10, 20, leopard.Code,
	}

	return JsonOutputHandler(http.StatusOK, sampleOutput)
}
