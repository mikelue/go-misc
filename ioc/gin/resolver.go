/*
You can customized behavior for "injected parameters".

The resolver would be used if:
  1. The type is not struct(or pointer to struct), and
  2. The field of struct is not tagged by Gin's specification of binding.

If the type of parameter is struct(or pointer to struct) and
none of the field is tagged by Gin's specification of binding,
the registered resolver would be tried.

Implements "ParamResolver"

You can use any value which implements "ParamResolver" to provide
customized behavior for parameter binding.

  type ResolveMyParam struct {}
  func (p *ResolveMyParam) CanResolve(sourceType reflect.Type) bool {
    return true;
  }
  func (p *ResolveMyParam) Resolve(context *gin.Context, sourceType reflect.Type) (interface{}, error) {
    return 0;
  }

  config.RegisterParamResolver(&ResolveMyParam{})

See "ParamAsFieldResolver" for resolving data in struct.

Implements "Resolvable"

The bound parameter(or output) could implements "Resolvable",
which would be used to construct value of parameter.

  type MyUser struct {}
  func (self *MyUser) Resolve(context *gin.Context) error {
  }

  func yourHandler(user *MyUser) OutputBody {
    return nil
  }

  builder.WrapToGinHandler(yourHandler)

See "ResolvableField" for resolving data in struct.

References:

  "MvcConfig.RegisterParamResolver(ParamResolver)" - Add resolver for parameters
  "MvcConfig.RegisterParamAsFieldResolver(ParamAsFieldResolver)" - Add resolver(as field) for parameters
*/

package gin

import (
	"reflect"
	"github.com/gin-gonic/gin"
	ur "github.com/mikelue/go-misc/utils/reflect"
)

// Constructs and resolves parameter fed to "MvcHandler"
type ParamResolver interface {
	// Checks whether or not the type of value could be resolved
	CanResolve(reflect.Type) bool
	// Constructs and initializes the value of parameter
	Resolve(*gin.Context, reflect.Type) (interface{}, error)
}

// Constructs and resolves parameter(for field of struct) fed to "MvcHandler"
type ParamAsFieldResolver interface {
	// Checks whether or not the value of field could be resolved
	CanResolve(*reflect.StructField) bool
	// Constructs and initializes the value of the struct's field
	Resolve(*gin.Context, *reflect.StructField) (interface{}, error)
}

// A type could implements this interface to customize resolving
type Resolvable interface {
	Resolve(*gin.Context) error
}

// A type could implements this interface to customize resolving
// when it is enclosed as field of a struct
type ResolvableField interface {
	ResolveField(*gin.Context, *reflect.StructField) error
}

type resolverController []ParamResolver
func (self resolverController) resolveBuilder(targetType reflect.Type) argvBuilder {
	var resolver ParamResolver
	for _, checkedResolver := range self {
		if checkedResolver.CanResolve(targetType) {
			resolver = checkedResolver
			break
		}
	}

	if resolver == nil {
		return nil
	}

	return func(c *gin.Context) (reflect.Value, error) {
		v, err := resolver.Resolve(c, targetType)

		var resolvedValue reflect.Value
		if v != nil {
			resolvedValue = reflect.ValueOf(v)
		} else {
			resolvedValue = reflect.Zero(targetType)
		}

		return resolvedValue, err
	}
}

var typeOfParamResolver reflect.Type = ur.TypeExtBuilder.NewByAny((*ParamResolver)(nil)).
	InterfaceType().AsType()
var typeOfParamAsFieldResolver reflect.Type = ur.TypeExtBuilder.NewByAny((*ParamAsFieldResolver)(nil)).
	InterfaceType().AsType()
var typeOfResolvable reflect.Type = ur.TypeExtBuilder.NewByAny((*Resolvable)(nil)).
	InterfaceType().AsType()
var typeOfResolvableField reflect.Type = ur.TypeExtBuilder.NewByAny((*ResolvableField)(nil)).
	InterfaceType().AsType()
