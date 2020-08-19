package gin

func ExampleNewMvcConfig_toBuilder() {
	builder := NewMvcConfig().ToBuilder()
	_ = builder
}
