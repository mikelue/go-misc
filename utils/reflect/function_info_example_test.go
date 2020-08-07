package reflect

import (
	"fmt"
)

func SampleFunc(a int) string {
	return fmt.Sprintf("A: %d", a)
}

func ExampleFuncInfo_inTypes() {
	typeExt := TypeExtBuilder.NewByAny(SampleFunc)

	fmt.Printf("In[%v]",
		typeExt.FuncInfo().InTypes()[0].Kind(),
	)
	// Output:
	// In[int]
}

func ExampleFuncInfo_outTypes() {
	typeExt := TypeExtBuilder.NewByAny(SampleFunc)

	fmt.Printf("Out[%v]",
		typeExt.FuncInfo().OutTypes()[0].Kind(),
	)
	// Output:
	// Out[string]
}
