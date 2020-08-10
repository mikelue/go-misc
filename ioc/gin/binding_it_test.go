package gin

import (
	"fmt"
	"reflect"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bindnig test over real Gin server", func() {
	var server *StoppableServer
	var ginEngine *gin.Engine

	mvcBuilder := NewMvcConfig().
		RegisterParamResolvers(turnipParamResolver(0)).
		RegisterErrorHandlers(onionsErrorHandler(0)).
		ToBuilder()

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		ginEngine = gin.Default()
		ginEngine.Use(
			func (c *gin.Context) {
				defer func() {
					if p := recover(); p != nil {
						fmt.Printf("\nRecover: %v\n", p)
					}
				}()

				c.Next()
			},
		)
	})
	JustBeforeEach(func() {
		server = NewStoppableServer(ginEngine)
		server.ListenAndServeAsync(testAddr)
	})
	AfterEach(func() {
		server.Shutdown()
		server = nil
		ginEngine = nil
	})

	Context("Bindings for all supported paramter types", func() {
		BeforeEach(func() {
			ginEngine.
				POST("/car/:driver_id", mvcBuilder.WrapToGinHandler(aCar))
			ginEngine.
				GET("/err-car", mvcBuilder.WrapToGinHandler(errCar))
		})

		It("[/car/{driver_id}] with (<Gin Binding>, <Resolvable>, <ParamResolver>)", func() {
			testedResp, err := resty.New().R().
				SetHeader("kc_key", "oVXIqnLcC9").
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]interface{} {
					"id": 323, "name": "Grea", "color": 3,
				}).
				SetResult(map[string]interface{}{}).
				Post(testUrl("/car/3881"))

			/**
			 * Asserts the response status
			 */
			Expect(err).To(Succeed())
			Expect(testedResp.StatusCode()).To(BeEquivalentTo(http.StatusOK))
			// :~)

			/**
			 * Asserts the response body
			 */
			testedResult := *(testedResp.Result().(*map[string]interface{}))
			Expect(testedResult["name"]).To(BeEquivalentTo("Grea"))
			Expect(testedResult["color"]).To(BeEquivalentTo(4))
			Expect(testedResult["kc_key"]).To(BeEquivalentTo("oVXIqnLcC9"))
			Expect(testedResult["driver_id"]).To(BeEquivalentTo(3881))
			Expect(testedResult["country_code"]).To(BeEquivalentTo("FbOHAtlyvp"))
			Expect(testedResult["vegatable"]).To(BeEquivalentTo("good-taste"))
			// :~)
		})

		It("[/err-car?passed=no] Error handler and status", func() {
			testedResp, err := resty.New().R().
				SetResult(map[string]interface{}{}).
				Get(testUrl("/err-car?passed=no"))

			/**
			 * Asserts the response status
			 */
			Expect(err).To(Succeed())
			Expect(testedResp.StatusCode()).To(BeEquivalentTo(http.StatusBadRequest))
			// :~)

			/**
			 * Asserts the response body
			 */
			Expect(testedResp.Header().Get("err_code")).To(BeEquivalentTo("981"))
			// :~)
		})

		It("[/err-car?passed=yes] status(Use Proxy)", func() {
			testedResp, err := resty.New().R().
				SetResult(map[string]interface{}{}).
				Get(testUrl("/err-car?passed=yes"))

			/**
			 * Asserts the response status
			 */
			Expect(err).To(Succeed())
			Expect(testedResp.StatusCode()).To(BeEquivalentTo(http.StatusUseProxy))
			// :~)
		})
	})
})

const testAddr = ":8080"
func testUrl(uri string) string {
	return fmt.Sprintf("http://localhost%s%s", testAddr, uri)
}

type iceland struct {
	code string
}

// Resolvable
func (v *iceland) Resolve(c *gin.Context) error {
	v.code = "FbOHAtlyvp"
	return nil
}

type turnip string

type turnipParamResolver int
func (turnipParamResolver) CanResolve(sourceType reflect.Type) bool {
	return sourceType.Name() == "turnip"
}
func (turnipParamResolver) Resolve(c *gin.Context, sourceType reflect.Type) (interface{}, error) {
	return turnip("good-taste"), nil
}

type onionsErrorHandler int
func (onionsErrorHandler) CanHandle(c *gin.Context, err error) bool {
	return err.Error() == "E-981"
}
func (onionsErrorHandler) HandleError(c *gin.Context, err error) error {
	c.Header("err_code", "981")
	return nil
}

func aCar(
	car *struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Color uint8 `json:"color"`
		KcKey string `header:"kc_key"`
		DriverId int `uri:"driver_id"`
	},
	country *iceland,
	vegatable turnip,
) OutputHandler {
	result := map[string]interface{} {
		"name": car.Name,
		"color": car.Color + 1,
		"kc_key": car.KcKey,
		"driver_id": car.DriverId,
		"country_code": country.code,
		"vegatable": vegatable,
	}

	return JsonOutputHandler(http.StatusOK, result)
}

func errCar(c *gin.Context) (int, error) {
	if c.Query("passed") == "yes" {
		return http.StatusUseProxy, nil
	}

	return http.StatusBadRequest, fmt.Errorf("E-981")
}
