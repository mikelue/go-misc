package dingo

import (
	"flamingo.me/dingo"

	fg "github.com/mikelue/go-misc/ioc/frangipani"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppContext", func() {
	testedContext := initTestedAppContext()

	It("Gets instance", func() {
		testedCar := testedContext.GetInstance(new(car)).(*car)

		Expect(testedCar.Engine.code).To(BeEquivalentTo("KWO-912"))
	})
	It("Gets environment object", func() {
		testeddProps := testedContext.Environment().Typed()

		Expect(testeddProps.GetInt("prop.key1")).To(BeEquivalentTo(20))
	})
})

type engine struct {
	code string
}
type car struct {
	Engine *engine `inject:""`
}

func initTestedAppContext() AppContext {
	sampleInjector, err := dingo.NewInjector()

	if err != nil {
		panic(err)
	}

	sampleInjector.Bind(new(engine)).
		ToInstance(&engine{ "KWO-912" })
	sampleInjector.Bind(new(car))
	sampleInjector.Bind(new(fg.Environment)).
		ToInstance(fg.EnvBuilder.NewByMap(
			map[string]interface{} {
				"prop.key1": 20,
				"prop.key2": 40,
			},
		))

	return AsAppContext(sampleInjector)
}
