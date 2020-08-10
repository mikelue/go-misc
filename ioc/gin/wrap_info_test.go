package gin

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wrap engine", func() {
	DescribeTable("getStructType",
		func(value interface{}, isStruct bool) {
			matcher := BeNil()
			testedResult := getStructType(reflect.TypeOf(value))

			if isStruct {
				matcher = Not(matcher)
			}

			Expect(testedResult).To(matcher)
		},
		Entry("Struct", hackberry{}, true),
		Entry("*Struct", &hackberry{}, true),
		Entry("**Struct", (&hackberry{}).getP(), true),
		Entry("int", 23, false),
		Entry("func() {}", func() {}, false),
	)

	Context("inTypes", func() {
		It("toBuilder(ParamResolver)", func() {
			sampleResolver := &sampleParamResolver{ true, 901 }
			sampleTypes := inTypes { reflect.TypeOf(sampleResolver) }
			values, err := toBuilderAndGetsValues(sampleTypes, sampleResolver)

			Expect(err).To(Succeed())
			Expect(values[0].Interface()).To(BeEquivalentTo(901))
		})
		It("toBuilder(Resolvable)", func() {
			sampleParam := parrot(20)
			sampleTypes := inTypes { reflect.TypeOf(&sampleParam) }

			values, err := toBuilderAndGetsValues(sampleTypes)

			Expect(err).To(Succeed())
			Expect(
				reflect.Indirect(values[0]).
					Interface(),
			).To(BeEquivalentTo(80))
		})
		It("toBuilder(pointer to struct)", func() {
			var sampleTypes inTypes = inTypes {
				reflect.TypeOf(&hackberryHeader{}),
				reflect.TypeOf(&hackberryBody{}),
				reflect.TypeOf(&hackberryUri{}),
			}

			// No panic
			values, err := toBuilderAndGetsValues(sampleTypes)
			Expect(err).To(Succeed())
			Expect(values).To(HaveLen(3))
			Expect(values[0].Interface().(*hackberryHeader).AcKey).To(BeEquivalentTo("OaVGJnls"))
			Expect(values[1].Interface().(*hackberryBody).Name).To(BeEquivalentTo("Ellis"))
			Expect(values[2].Interface().(*hackberryUri).ObjectId).To(BeEquivalentTo(465))
		})
		It("toBuilder(struct)", func() {
			var sampleTypes inTypes = inTypes {
				reflect.TypeOf(hackberryHeader{}),
				reflect.TypeOf(hackberryBody{}),
				reflect.TypeOf(hackberryUri{}),
			}

			// No panic
			values, err := toBuilderAndGetsValues(sampleTypes)
			Expect(err).To(Succeed())
			Expect(values).To(HaveLen(3))
			Expect(values[0].Interface().(hackberryHeader).AcKey).To(BeEquivalentTo("OaVGJnls"))
			Expect(values[1].Interface().(hackberryBody).Name).To(BeEquivalentTo("Ellis"))
			Expect(values[2].Interface().(hackberryUri).ObjectId).To(BeEquivalentTo(465))
		})
		It("toBuilder(*gin.Context)", func() {
			var sampleTypes inTypes = inTypes { typeOfGinContext }

			// No panic
			values, err := toBuilderAndGetsValues(sampleTypes)
			Expect(err).To(Succeed())
			Expect(values).To(HaveLen(1))
			Expect(values[0].Interface().(*gin.Context)).ToNot(BeNil())
		})

		DescribeTable("toBuilder(error of un-supported types)",
			func(v interface{}, errorPattern string) {
				sampleTypes := inTypes { reflect.TypeOf(v) }

				// With panic
				Expect(
					func() { sampleTypes.toBuilder(make(resolverController, 0)) },
				).To(PanicWith(MatchError(MatchRegexp(errorPattern))))
			},
			Entry("Un-supported type", 32, `Args\[0\].*`),
			Entry("Struct is not supported by Gin's binding", &hackberry{}, "Unable to find supported.*Gin"),
		)
	})

	Context("outTypes", func() {
		It("toCallbacks", func() {
			sampleTypes := outTypes {
				reflect.TypeOf(200),
				reflect.TypeOf(fmt.Errorf("sample-error")),
				reflect.TypeOf((*OutputHandler)(nil)).Elem(),
			}

			Expect(sampleTypes.toCallbacks()).To(HaveLen(3))
		})
		It("toCallbacks(error)", func() {
			sampleTypes := outTypes {
				reflect.TypeOf(200),
				reflect.TypeOf("This is not supported"),
			}

			// With panic
			Expect(
				func() { sampleTypes.toCallbacks() },
			).To(PanicWith(MatchError(MatchRegexp(`value\[1].*string`))))
		})
	})

	Context("Output callbacks", func() {
		var resp *httptest.ResponseRecorder
		var sampleContext *gin.Context

		BeforeEach(func() {
			resp = httptest.NewRecorder()
			sampleContext, _ = gin.CreateTestContext(resp)
		})

		DescribeTable("errorOutputCallback",
			func(err error) {
				matcher := Succeed()
				if err != nil {
					matcher = MatchError(err)
				}

				Expect(errorOutputCallback(nil, err)).To(matcher)
			},
			Entry("Viable", fmt.Errorf("sample error")),
			Entry("Nil", nil),
		)

		It("statusOutputCallback", func() {
			statusOutputCallback(sampleContext, 201)

			Expect(sampleContext.Writer.Status()).To(BeEquivalentTo(201))
		})

		It("outputBodyCallback", func() {
			var sampleOutputHandler OutputHandlerFunc = func(context *gin.Context) error {
				context.Status(303)
				return nil
			}
			outputHandlerCallback(sampleContext, sampleOutputHandler)

			Expect(sampleContext.Writer.Status()).To(BeEquivalentTo(303))
		})
	})

	Context("resolveByGinBinding", func() {
		var sampleContext *gin.Context

		BeforeEach(func() {
			resp := httptest.NewRecorder()
			sampleContext, _ = gin.CreateTestContext(resp)
			sampleContext.Request = httptest.NewRequest("POST", "/test-1", strings.NewReader(`{ "id": 24, "name": "Animal Keeper" }`))
			sampleContext.Params = gin.Params{ { Key: "object_id", Value: "761" } }
			sampleContext.Request.Header["Content-Type"] = []string{ "application/json" }
			sampleContext.Request.Header["Adkey"] = []string{ "Yn1CvIXa" }
			sampleContext.Request.Header["Ackey"] = []string{ "nz4Z8nKy" }
		})

		It("shouldBindCallback", func() {
			testedValue, err := resolveByGinBinding(reflect.TypeOf(hackberryBody{}))(sampleContext)
			Expect(err).To(Succeed())

			testedDedicateValue := testedValue.Interface().(*hackberryBody)

			Expect(testedDedicateValue.Id).To(BeEquivalentTo(24))
			Expect(testedDedicateValue.Name).To(BeEquivalentTo("Animal Keeper"))
		})

		It("shouldBindUriCallback", func() {
			testedValue, err := resolveByGinBinding(reflect.TypeOf(hackberryUri{}))(sampleContext)
			Expect(err).To(Succeed())

			testedDedicateValue := testedValue.Interface().(*hackberryUri)
			Expect(testedDedicateValue.ObjectId).To(BeEquivalentTo(761))
		})

		It("shouldBindHeaderCallback", func() {
			testedValue, err := resolveByGinBinding(reflect.TypeOf(hackberryHeader{}))(sampleContext)
			Expect(err).To(Succeed())

			testedDedicateValue := testedValue.Interface().(*hackberryHeader)
			Expect(testedDedicateValue.AdKey).To(BeEquivalentTo("Yn1CvIXa"))
			Expect(testedDedicateValue.AcKey).To(BeEquivalentTo("nz4Z8nKy"))
		})
	})

	DescribeTable("isResolvable",
		func(sampleValueForType interface{}, expectedResult bool) {
			Expect(isResolvable(reflect.TypeOf(sampleValueForType))).
				To(BeEquivalentTo(expectedResult))
		},
		Entry("Type is resovlable", (*parrot)(nil), true),
		Entry("Pointer to the type is resovlable", (parrot)(0), true),
		Entry("The type is not resolvable", 0, false),
	)

	Context("resovlableBuilder", func() {
		DescribeTable("Resolving passed",
			func(sampleValueForType interface{}) {
				sampleType := reflect.TypeOf(sampleValueForType)
				testedBuilder := resovlableBuilder(sampleType)

				context, _ := newContext();
				value, err := testedBuilder(context);

				Expect(err).To(Succeed())
				Expect(*(value.Interface().(*parrot))).To(BeEquivalentTo(80))
			},
			Entry("Type is resovlable", (*parrot)(nil)),
			Entry("Pointer to the type is resovlable", (parrot)(0)),
		)
		It("Resolving gives error", func() {
			testedBuilder := resovlableBuilder(reflect.TypeOf((*parrot)(nil)))

			context, _ := newContext();
			context.Params = append(context.Params, gin.Param{ Key: "k1", Value: "981" })
			_, err := testedBuilder(context);

			Expect(err).To(MatchError("Sample-Error"))
		})
	})

	Context("bindTypeFlag", func() {
		var sampleType bindTypeFlag = 0x01
		structType := reflect.TypeOf(
			struct {
				P1 int `t1:"20" t2:"40"`
			} {},
		)

		DescribeTable("Mask with various situation",
			func(v int, fieldName string, tagNames []string, expected int) {
				field, _ := structType.FieldByName(fieldName)
				Expect(sampleType.matchOr(bindTypeFlag(v), field.Tag, tagNames...)).
					To(BeEquivalentTo(expected))
			},
			Entry("Mask matched", 0x00, "P1", []string{ "t1" }, 0x01),
			Entry("Mask matched", 0x00, "P1", []string{ "t3" }, 0x00),
			Entry("Mask matched", 0x03, "P1", []string{ "t1" }, 0x03),
			Entry("Mask matched", 0x03, "P1", []string{ "t3" }, 0x03),
		)
	})
})

func toBuilderAndGetsValues(types inTypes, resolvers ...ParamResolver) ([]reflect.Value, error) {
	resp := httptest.NewRecorder()
	sampleContext, _ := gin.CreateTestContext(resp)
	sampleContext.Request = httptest.NewRequest("POST", "/test-2", strings.NewReader(`{ "id": 93, "name": "Ellis" }`))
	sampleContext.Request.Header["Content-Type"] = []string{ "application/json" }
	sampleContext.Request.Header["Adkey"] = []string{ "kEFU2Ss7" }
	sampleContext.Request.Header["Ackey"] = []string{ "OaVGJnls" }
	sampleContext.Params = gin.Params{ { Key: "object_id", Value: "465" } }

	return types.toBuilder(resolverController(resolvers))(sampleContext)
}

func sampleFunc1(a int8, b int64, c string) (int16, error) {
	return 0, nil
}

func sampleOutputFunc(context *gin.Context) {}

type parrot int
func (self *parrot) Resolve(context *gin.Context) error {
	if _, ok := context.Params.Get("k1"); ok {
		return fmt.Errorf("%v", "Sample-Error")
	}

	*self = 80
	return nil
}

type hackberry struct {}
func (self *hackberry) getP() **hackberry {
	return &self
}

type hackberryHeader struct {
	AdKey string `header:"Adkey"`
	AcKey string `header:"Ackey"`
}
type hackberryBody struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
type hackberryUri struct {
	ObjectId int `uri:"object_id"`
}
