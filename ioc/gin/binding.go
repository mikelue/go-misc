/*
MvcHandler

A MVC handler can be any function with ruled types of parameter and defined returned types.

  type MvcHandler interface{}

You can define handler of supported:

  func(req *http.Request, params *gin.Params) OutputHandler {
    return TextOutputHandler("Hello World")
  }

  func(
    data *struct {
      // Binding data from form(or query string)
      Name string `form:"name" default:"anonymous"`
      // Binding data from uri
      Age int `uri:"id" default:"0"`
      // Binding data from header
      SessionId string `header:"session_id"`
      // Binding data from JSON
      Weight int `json:"weight"`
    },
  ) OutputHandler {
    return TextOutputHandler("Hello World")
  }

WrapToGinHandler

After you define the MVC handler, you could use "MvcBuilder.WrapToGinHandler()" to
convert your handler to "gin.HandlerFunc".

  mvcBuilder := NewMvcConfig().ToBuilder()
  engine.Get("/your-web-service", mvcBuilder.WrapToGinHandler(your_mvc_handler))

Parameters of Handler

"<struct>" - See parameter tags for automatic binding
 This type of value woule be checked by ogin.ConformAndValidateStruct automatically.

"*gin.Context" - The context object of current request

TODO-Parameters of Handler

"json.Unmarshaler" - If the type of value is json.Unmarshaler, use the UnmarshalJSON([]byte) function of the value
 This type of value woule be checked by ogin.ConformAndValidateStruct automatically.

"gin.ResponseWriter" - See "gin.ResponseWriter"

"gin.Params" - See "gin.Params"

"*http.Request" - See "http.Request"

"http.ResponseWriter" - See "http.ResponseWriter"

"*url.URL" - See "url.URL"

"*multipart.Reader" - See "multipart.Reader"; Once you use *multipart.Form, the reader would reach EOF.

"*multipart.Form" - See "multipart.Form"

"*validator.Validate" - See go-playground/validator.v10

Tagging Struct

There are various definition of tags could be used on struct:

  type MyData struct {
    // Binding data from form(or query string)
    Name string `form:"name"`
    // Binding data from uri
    Age int `uri:"id"`
    // Binding data from header
    SessionId string `header:"session_id"`
    // Binding data from JSON
    Weight int `json:"weight"`
  }

Form(Query string), header, or JSON body

  form:"field_1" - Use query parameter param_name_1 as binding value
  form:"field_2" - Must be bool type, used to indicate whether or not has viable value for this parameter
  header:"header_value_1" - Use the value of URI parameter pm_1 as binding value
  header:"header_value_2" - Use the form value of in_1 as binding value
  uri:"uri_v_1" - Must be bool type, used to indicate whether or not has viable value for this parameter
  uri:"uri_v_2" - Use the header value of Content-Type as binding value
  json:"v1" - Must be bool type, used to indicate whether or not has viable value for this parameter
  json:"v2" - Must be bool type, used to indicate whether or not has viable value for this parameter

TODO-Default Value

  default:"20" - Gives value 20 if the value of binding is empty
  default:"[20,40,30]" - Gives value [20, 40, 30](as array, no space)if the value of binding is empty

By default, if the value of binding is existing, the framework would use the default value of binding type.

Data Validation

The Gin framework uses "go-playground/validator"(v10) as default validator.

See Also:
Documentations of Gin's Binding: https://github.com/gin-gonic/gin#model-binding-and-validation
*/
package gin

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	ur "github.com/mikelue/go-misc/utils/reflect"
)

// Constructs MVC configuration with default values
// 	EnableValidator(bool) - true
func NewMvcConfig() *MvcConfig {
	return &MvcConfig{
		paramResolvers: make([]ParamResolver, 0, 4),
		paramAsFieldResolvers: make([]ParamAsFieldResolver, 0, 4),
		errorController: make([]ErrorHandler, 0, 2),
	}
}

// The configuration used by "MvcBuilder"
type MvcConfig struct {
	paramResolvers []ParamResolver
	paramAsFieldResolvers []ParamAsFieldResolver
	errorController errorController
}

// Registers multiple resolvers
func (self *MvcConfig) RegisterParamResolvers(paramResolvers ...ParamResolver) *MvcConfig {
	self.paramResolvers = append(self.paramResolvers, paramResolvers...)
	return self
}
// Registers multiple resolvers of field
func (self *MvcConfig) RegisterParamAsFieldResolvers(paramAsFieldResolvers ...ParamAsFieldResolver) *MvcConfig {
	self.paramAsFieldResolvers = append(self.paramAsFieldResolvers, paramAsFieldResolvers...)
	return self
}
// Registers multiple handlers for error
//
// By default, any unhandled error would be output as JSON and 500(Internal server error) status
func (self *MvcConfig) RegisterErrorHandlers(errHandlers ...ErrorHandler) *MvcConfig {
	self.errorController = append(self.errorController, errHandlers...)
	return self
}

// Gets the instance of "MvcBuilder"
func (self *MvcConfig) ToBuilder() MvcBuilder {
	clonedConfig := *self
	clonedConfig.errorController = append(
		self.errorController, defaultErrorHandler(0),
	)

	return &mvcBuilderImpl{
		config: &clonedConfig,
	}
}

// As alias for "interface{}"(GoLang Sucks)
type MvcHandler = interface{}

// This build warp MVC handler to gin handler(as "func(c *gin.Context)")
type MvcBuilder interface {
	// The wrapper method
	WrapToGinHandler(MvcHandler) func(c *gin.Context)
}

type mvcBuilderImpl struct {
	config *MvcConfig
}

func (self *mvcBuilderImpl) WrapToGinHandler(mvcHandler MvcHandler) func(c *gin.Context) {
	/**
	 * In arguments, Out variables and function value for performing calling
	 */
	funcInfo := ur.TypeExtBuilder.NewByAny(mvcHandler).FuncInfo()
	argsBuilder := inTypes(funcInfo.InAsTypes()).toBuilder(self.config.paramResolvers)
	outCallbacks := outTypes(funcInfo.OutAsTypes()).toCallbacks()
	funcValue := reflect.ValueOf(mvcHandler)
	// :~)

	return func(c *gin.Context) {
		args, err := argsBuilder(c)

		if err != nil {
			self.config.errorController.handle(c, err)
			return
		}

		returnedValues := funcValue.Call(args)

		for i, value := range returnedValues {
			if outErr := outCallbacks[i](c, value.Interface()); outErr != nil {
				self.config.errorController.handle(c, outErr)
				return
			}
		}
	}
}

type defaultErrorHandler int
func (defaultErrorHandler) CanHandle(context *gin.Context, err error) bool {
	return true
}
func (defaultErrorHandler) HandleError(context *gin.Context, err error) error {
	context.JSON(http.StatusInternalServerError, err)
	return nil
}
