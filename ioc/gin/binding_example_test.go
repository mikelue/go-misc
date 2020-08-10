package gin

import (
	"net/http"
)

func ExampleNewMvcConfig_toBuilder() {
	builder := NewMvcConfig().ToBuilder()
	_ = builder
}

func ExampleMvcBuilder_wrapToGinHandler() {
	builder := NewMvcConfig().ToBuilder()
	handler := builder.WrapToGinHandler(sampleHandler)

	// ginEngine.Post("/some-mammal", handler)
	_ = handler
}

type leopard struct {
	Id int `json:"id"`
	Color uint8 `json:"color"`
	Height int `json:"height"`
	Weight int `json:"weight"`
}

func sampleHandler() OutputHandler {
	return JsonOutputHandler(http.StatusOK, `[10, 20, 30]`)
}
