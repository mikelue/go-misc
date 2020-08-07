/*
TrueValue

As the instance of "reflect.Value" for "true"

  someValue := TrueValue

FalseValue

  someValue := FalseValue

As the instance of "reflect.Value" for "false"
*/
package types

import (
	"reflect"
)

var (
	/**
	 * Boolean values
	 */
	TrueValue  = reflect.ValueOf(true)
	FalseValue = reflect.ValueOf(false)
	// :~)
)
