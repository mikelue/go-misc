/*
OutputHandler

You can provides various combination of returned value for "MvcHandler":

  OutputHandler - Single value of output body
  (OutputHandler, error) - OutputHandler nad error
  (OutputHandler, int) - response status and OutputHandler
  (OutputHandler, int, error) - response status, OutputHandler, and error

Implements "OutputHandler"

You can implements "OutputHandler" of returned type:

  type MyCar struct {}
  func (self *MyUser) Output(context *gin.Context) error {
    context.JSON(200, "[1, 11, 21]")
    return nil
  }

  func yourHandler() *MyCar {
    return nil
  }

see: "JsonOutputHandler", "TextOutputHandler", "XmlOutputHandler", etc.

TODO-Return value by other types

"json.Marshaler" - If the type of returned value is json.Marshaler, use JsonOutputHandler() as output type

"string" - If the type of returned value is string, use TextOutputHandler() as output type

"fmt.Stringer" - As same as string
*/
package gin

import (
	"mime"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ur "github.com/mikelue/go-misc/utils/reflect"
)

// Main interface for generating response
type OutputHandler interface {
	Output(*gin.Context) error
}

// Functional type of "OutputHandler"
type OutputHandlerFunc func(*gin.Context) error

// As implementation of "OutputHandler"
func (f OutputHandlerFunc) Output(context *gin.Context) error {
	return f(context)
}

// With HTTP header of "Accepted", this function builds "OutputHandler" by the supported MIME type:
//
//  application/json, application/xml, text/xml, text/plain,
//	application/x-protobuf, application/x-yaml
func AutoDetectOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		outputHandler := builderByAccept(context)(code, v)
		return outputHandler.Output(context)
	})
}

type outputHandlerBuilder func(int, interface{}) OutputHandler

func builderByAccept(context *gin.Context) outputHandlerBuilder {
	for _, acceptContentType := range context.Request.Header.Values("Accept") {
		mediaType, _, err := mime.ParseMediaType(acceptContentType)
		if err != nil {
			continue
		}

		switch mediaType {
			case binding.MIMEJSON:
				return JsonOutputHandler
			case binding.MIMEXML, binding.MIMEXML2:
				return XmlOutputHandler
			case binding.MIMEPlain:
				return TextOutputHandler
			case binding.MIMEPROTOBUF:
				return ProtoBufOutputHandler
			case binding.MIMEYAML:
				return YamlOutputHandler
		}
	}

	return nil
}

// Uses "(*gin.Context).JSON(http.StatusOK, v)" to perform response
func JsonOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.JSON(code, v)
		return nil
	})
}

// Uses "(*gin.Context).String(code, v)" to perform response
func TextOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.String(code, "%s", v)
		return nil
	})
}

// Uses "(*gin.Context).HTML(code, name, v)" to perform response
func HtmlOutputHandler(code int, name string, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.HTML(code, name, v)
		return nil
	})
}

// Uses "(*gin.Context).XML(code, v)" to perform response
func XmlOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.XML(code, v)
		return nil
	})
}

// Uses "(*gin.Context).YAML(code, v)" to perform response
func YamlOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.YAML(code, v)
		return nil
	})
}

// Uses "(*gin.Context).ProtoBuf(code, v)" to perform response
func ProtoBufOutputHandler(code int, v interface{}) OutputHandler {
	return OutputHandlerFunc(func(context *gin.Context) error {
		context.ProtoBuf(code, v)
		return nil
	})
}

var outputHandlerType reflect.Type = ur.TypeExtBuilder.NewByAny((*OutputHandler)(nil)).
	InterfaceType().AsType()
