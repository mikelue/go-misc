package gin

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	ur "github.com/mikelue/go-misc/utils/reflect"
	tr "github.com/mikelue/go-misc/utils/reflect/types"
)

// Constructs parameters of "injected" at warm-up time
type inTypes []reflect.Type
// As list of builders for arguments
func (self inTypes) toBuilder(externalResolver resolverController) argsBuilder {
	argsBuilderAsFuncs := make([]argvBuilder, 0, len(self))

	/**
	 * Collects building functions for arguments
	 */
	for i, inType := range self {
		if inType.AssignableTo(typeOfGinContext) {
			argsBuilderAsFuncs = append(argsBuilderAsFuncs, ginContextBuilder)
			continue
		}

		/**
		 * Supports "ParamResolver"
		 */
		if externalArgvBuilder := externalResolver.resolveBuilder(inType); externalArgvBuilder != nil {
			argsBuilderAsFuncs = append(argsBuilderAsFuncs, externalArgvBuilder)
			continue
		}
		// :~)

		/**
		 * Supports "Resolvable"
		 */
		if isResolvable(inType) {
			argsBuilderAsFuncs = append(argsBuilderAsFuncs, resovlableBuilder(inType))
			continue
		}
		// :~)

		/**
		 * Build-in resolving by struct
		 */
		structType := getStructType(inType)
		if structType != nil {
			argsBuilderAsFuncs = append(argsBuilderAsFuncs, resolveByGinBinding(structType))
			continue
		}
		// :~)

		panic(fmt.Errorf("Args[%d] needs to be struct or \"Resolvable\" or registering of \"ParamResolver\". But got: \"%v\".", i, inType))
	}
	// :~)

	return func(context *gin.Context) ([]reflect.Value, error) {
		values := make([]reflect.Value, 0, len(self))

		for i, builderFunc := range argsBuilderAsFuncs {
			value, err := builderFunc(context)

			if err != nil {
				return nil, err
			}

			/**
			 * If the desired type is the concrete value of pointer,
			 * converts it to concrete element.
			 */
			valueType := value.Type()
			if !valueType.AssignableTo(self[i]) &&
				valueType.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			// :~)

			values = append(values, value)
		}

		return values, nil
	}
}

// Constructs output of "MvcHandler" at warm-up time
type outTypes []reflect.Type
// As list of builders for output returned values of "MvcHandler"
func (self outTypes) toCallbacks() []outputCallback {
	builders := make([]outputCallback, 0, len(self))

	for i, outType := range self {
		if outType.Implements(outputHandlerType) {
			builders = append(builders, outputHandlerCallback)
		} else if outType.Implements(tr.ErrorType) {
			builders = append(builders, errorOutputCallback)
		} else if outType.Kind() == reflect.Int {
			builders = append(builders, statusOutputCallback)
		} else {
			panic(fmt.Errorf(
				"Unsupported type of returned value[%d]: %v. Supporting types: [int, OutputHandler, error]",
				i, outType,
			))
		}
	}

	return builders
}

var (
	typeOfGinContext = reflect.TypeOf((*gin.Context)(nil))
)

type argvBuilder func(context *gin.Context) (reflect.Value, error)
type argsBuilder func(context *gin.Context) ([]reflect.Value, error)
type outputCallback func(context *gin.Context, returnedValue interface{}) error

// Checks if the type of value is "Resolvable" or
// the pointer to the value is "Resolvable"
func isResolvable(targetType reflect.Type) bool {
	if targetType.Implements(typeOfResolvable) {
		return true
	}

	if reflect.PtrTo(targetType).Implements(typeOfResolvable) {
		return true
	}

	return false
}

// Supports two shapes of type:
//
//  1. pointer is Resolvable(one-level)
//  2. pointer of value is Resolvable(one-level)
func resovlableBuilder(targetType reflect.Type) argvBuilder {
	if targetType.Implements(typeOfResolvable) && targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	return func(context *gin.Context) (reflect.Value, error) {
		newValue := reflect.New(targetType)
		err := newValue.Interface().(Resolvable).Resolve(context)
		return newValue, err
	}
}

// As "argvBuilder"
func ginContextBuilder(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context), nil
}

func getStructType(valueType reflect.Type) reflect.Type {
	finalType := ur.TypeExtBuilder.NewByType(valueType).
		RecursiveIndirect().AsType()

	if finalType.Kind() != reflect.Struct {
		return nil
	}

	return finalType
}

func resolveByGinBinding(structType reflect.Type) argvBuilder {
	/**
	 * Figures out the properties(supported by Gin binding) of tag
	 */
	var typeFlags bindTypeFlag = 0
	for i := 0; i < structType.NumField(); i++ {
		tag := structType.Field(i).Tag

		typeFlags = tfUri.matchOr(typeFlags, tag, "uri")
		typeFlags = tfHeader.matchOr(typeFlags, tag, "header")
		typeFlags = tfBody.matchOr(typeFlags, tag, "json", "xml", "form")
	}
	// :~)

	if typeFlags == 0 {
		panic(fmt.Errorf("Unable to find supported tag of Gin: %v", structType))
	}

	bindingCallbacks := make([]bindingCallback, 0, 1)

	/**
	 * Performs binding by tagging information
	 *
	 * WARNING: Because "gin.Context.ShouldBind" would override the struct,
	 *   the "shouldBindCallback" **MUST BE** the first callback.
	 */
	if (typeFlags & tfBody) > 0 {
		bindingCallbacks = append(bindingCallbacks, shouldBindCallback)
	}
	if (typeFlags & tfUri) > 0 {
		bindingCallbacks = append(bindingCallbacks, shouldBindUriCallback)
	}
	if (typeFlags & tfHeader) > 0 {
		bindingCallbacks = append(bindingCallbacks, shouldBindHeaderCallback)
	}
	// :~)

	return func(context *gin.Context) (reflect.Value, error) {
		newValueOfStruct := reflect.New(structType)

		for _, callback := range bindingCallbacks {
			if err := callback(context, newValueOfStruct.Interface()); err != nil {
				return reflect.Value{}, err
			}
		}

		return newValueOfStruct, nil
	}
}

type bindTypeFlag int
func (self bindTypeFlag) matchOr(v bindTypeFlag, tag reflect.StructTag, tagNames ...string) bindTypeFlag {
	/**
	 * The flags has already set
	 */
	if (v & self) > 0 {
		return v
	}
	// :~)

	ok := false
	for _, tagName := range tagNames {
		_, ok = tag.Lookup(tagName)
		if ok {
			break
		}
	}

	/**
	 * Found corresponding tag, mask the value
	 */
	if ok {
		return v | self
	}
	// :~)

	return v
}

const (
	tfUri bindTypeFlag = 0x01
	tfHeader bindTypeFlag = 0x02
	tfBody bindTypeFlag = 0x04
)

type bindingCallback func(*gin.Context, interface{}) error

func shouldBindCallback(context *gin.Context, value interface{}) error {
	return context.ShouldBind(value)
}
func shouldBindUriCallback(context *gin.Context, value interface{}) error {
	return context.ShouldBindUri(value)
}
func shouldBindHeaderCallback(context *gin.Context, value interface{}) error {
	return context.ShouldBindHeader(value)
}

func errorOutputCallback(context *gin.Context, v interface{}) error {
	if v == nil {
		return nil
	}

	return v.(error)
}

func statusOutputCallback(context *gin.Context, v interface{}) error {
	status, ok := v.(int)

	if !ok {
		return fmt.Errorf("Expect status(int), but got: %v", reflect.TypeOf(v))
	}

	context.Status(status)
	return nil
}
func outputHandlerCallback(context *gin.Context, v interface{}) error {
	outputBody, ok := v.(OutputHandler)

	if !ok {
		return fmt.Errorf("Expect OutputHandler, but got: %v", reflect.TypeOf(v))
	}

	return outputBody.Output(context)
}
