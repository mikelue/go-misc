package frangipani

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
	bs "github.com/inhies/go-bytesize"
)

// Method space used to construct new instance of PropertyResolver
var PropertyResolverBuilder IPropertyResolverBuilder = IPropertyResolverBuilder(0)

type IPropertyResolverBuilder int
// Constructs new resolvers by a map(elements are cloned shallowly)
func (IPropertyResolverBuilder) NewByMap(props map[string]interface{}) PropertyResolver {
	newProps := make(map[string]interface{}, len(props))

	for k, v := range props {
		newProps[k] = v
	}

	return mapBasedPropertyResolver(newProps)
}

// Resolves value of properties
type PropertyResolver interface {
	// Interface space for retrieving typed values of properties
	Typed() TypedR
	// Interface space for retrieving required and typed values of properties
	RequiredTyped() RequiredTypedR

	// Checks whether or not the property is existing
	ContainsProperty(string) bool
	// Gets the property as string(or empty)
	GetProperty(string) string

	// Gets the property as string
	//
	// The "error" is viable if "PropertyResolver.ContainsProperty()" returns false.
	GetRequiredProperty(string) (string, error)
}

// Defines the getting of property for specific types.
//
// These methods is cloned from "*viper.Viper".
//
// See "PropertyResolver.Typed()"
type TypedR interface {
	Get(string) interface{}
	GetBool(string) bool
	GetDuration(string) time.Duration
	GetFloat64(string) float64
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetIntSlice(string) []int
	GetString(string) string
	GetStringMap(string) map[string]interface{}
	GetStringMapString(string) map[string]string
	GetStringMapStringSlice(string) map[string][]string
	GetStringSlice(string) []string
	GetTime(string) time.Time
	GetUint(string) uint
	GetUint32(string) uint32
	GetUint64(string) uint64

	// See: github.com/inhies/go-bytesize
	GetByteSize(string) bs.ByteSize
}

// Defines the getting of property for specific types
// (with error if the property is not existing).
//
// These methods is cloned from "*viper.Viper"
//
// See "PropertyResolver.RequiredTyped()"
type RequiredTypedR interface {
	Get(string) (interface{}, error)
	GetBool(string) (bool, error)
	GetDuration(string) (time.Duration, error)
	GetFloat64(string) (float64, error)
	GetInt(string) (int, error)
	GetInt32(string) (int32, error)
	GetInt64(string) (int64, error)
	GetIntSlice(string) ([]int, error)
	GetString(string) (string, error)
	GetStringMap(string) (map[string]interface{}, error)
	GetStringMapString(string) (map[string]string, error)
	GetStringMapStringSlice(string) (map[string][]string, error)
	GetStringSlice(string) ([]string, error)
	GetTime(string) (time.Time, error)
	GetUint(string) (uint, error)
	GetUint32(string) (uint32, error)
	GetUint64(string) (uint64, error)

	// See: github.com/inhies/go-bytesize
	GetByteSize(string) (bs.ByteSize, error)
}

type mapBasedPropertyResolver map[string]interface{}

func (self mapBasedPropertyResolver) Typed() TypedR {
	return typedRImpl(self)
}
func (self mapBasedPropertyResolver) RequiredTyped() RequiredTypedR {
	return requiredTypedRImpl(self)
}
func (self mapBasedPropertyResolver) ContainsProperty(name string) bool {
	_, ok := self[name]
	return ok
}
func (self mapBasedPropertyResolver) GetProperty(name string) string {
	return self.Typed().GetString(name)
}
func (self mapBasedPropertyResolver) GetRequiredProperty(name string) (string, error) {
	return self.RequiredTyped().GetString(name)
}

type typedRImpl mapBasedPropertyResolver
func (self typedRImpl) Get(name string) interface{} {
	return self[name]
}
func (self typedRImpl) GetBool(name string) bool {
	return cast.ToBool(self.Get(name))
}
func (self typedRImpl) GetDuration(name string) time.Duration {
	return cast.ToDuration(self.Get(name))
}
func (self typedRImpl) GetFloat64(name string) float64 {
	return cast.ToFloat64(self.Get(name))
}
func (self typedRImpl) GetInt(name string) int {
	return cast.ToInt(self.Get(name))
}
func (self typedRImpl) GetInt32(name string) int32 {
	return cast.ToInt32(self.Get(name))
}
func (self typedRImpl) GetInt64(name string) int64 {
	return cast.ToInt64(self.Get(name))
}
func (self typedRImpl) GetIntSlice(name string) []int {
	return cast.ToIntSlice(self.Get(name))
}
func (self typedRImpl) GetString(name string) string {
	return cast.ToString(self.Get(name))
}
func (self typedRImpl) GetStringMap(name string) map[string]interface{} {
	return cast.ToStringMap(self.Get(name))
}
func (self typedRImpl) GetStringMapString(name string) map[string]string {
	return cast.ToStringMapString(self.Get(name))
}
func (self typedRImpl) GetStringMapStringSlice(name string) map[string][]string {
	return cast.ToStringMapStringSlice(self.Get(name))
}
func (self typedRImpl) GetStringSlice(name string) []string {
	return cast.ToStringSlice(self.Get(name))
}
func (self typedRImpl) GetTime(name string) time.Time {
	return cast.ToTime(self.Get(name))
}
func (self typedRImpl) GetUint(name string) uint {
	return cast.ToUint(self.Get(name))
}
func (self typedRImpl) GetUint32(name string) uint32 {
	return cast.ToUint32(self.Get(name))
}
func (self typedRImpl) GetUint64(name string) uint64 {
	return cast.ToUint64(self.Get(name))
}
func (self typedRImpl) GetByteSize(name string) bs.ByteSize {
	v, err := bs.Parse(self.GetString(name))
	if err != nil {
		return 0
	}

	return v
}

type requiredTypedRImpl mapBasedPropertyResolver
func (self requiredTypedRImpl) Get(name string) (interface{}, error) {
	v, ok := self[name]

	if !ok {
		return nil, fmt.Errorf("Property[%s] is not existing", name)
	}

	return v, nil
}
func (self requiredTypedRImpl) GetBool(name string) (bool, error) {
	v, err := self.Get(name)
	if err != nil {
		return false, err
	}

	return cast.ToBoolE(v)
}
func (self requiredTypedRImpl) GetDuration(name string) (time.Duration, error) {
	v, err := self.Get(name)
	if err != nil {
		return time.Duration(0), err
	}

	return cast.ToDurationE(v)
}
func (self requiredTypedRImpl) GetFloat64(name string) (float64, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToFloat64E(v)
}
func (self requiredTypedRImpl) GetInt(name string) (int, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToIntE(v)
}
func (self requiredTypedRImpl) GetInt32(name string) (int32, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt32E(v)
}
func (self requiredTypedRImpl) GetInt64(name string) (int64, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt64E(v)
}
func (self requiredTypedRImpl) GetIntSlice(name string) ([]int, error) {
	v, err := self.Get(name)
	if err != nil {
		return []int{}, err
	}

	return cast.ToIntSliceE(v)
}
func (self requiredTypedRImpl) GetString(name string) (string, error) {
	v, err := self.Get(name)
	if err != nil {
		return "", err
	}

	return cast.ToStringE(v)
}
func (self requiredTypedRImpl) GetStringMap(name string) (map[string]interface{}, error) {
	v, err := self.Get(name)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return cast.ToStringMapE(v)
}
func (self requiredTypedRImpl) GetStringMapString(name string) (map[string]string, error) {
	v, err := self.Get(name)
	if err != nil {
		return map[string]string{}, err
	}

	return cast.ToStringMapStringE(v)
}
func (self requiredTypedRImpl) GetStringMapStringSlice(name string) (map[string][]string, error) {
	v, err := self.Get(name)
	if err != nil {
		return map[string][]string{}, err
	}

	return cast.ToStringMapStringSliceE(v)
}
func (self requiredTypedRImpl) GetStringSlice(name string) ([]string, error) {
	v, err := self.Get(name)
	if err != nil {
		return []string{}, err
	}

	return cast.ToStringSliceE(v)
}
func (self requiredTypedRImpl) GetTime(name string) (time.Time, error) {
	v, err := self.Get(name)
	if err != nil {
		return time.Unix(0, 0), err
	}

	return cast.ToTimeE(v)
}
func (self requiredTypedRImpl) GetUint(name string) (uint, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUintE(v)
}
func (self requiredTypedRImpl) GetUint32(name string) (uint32, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint32E(v)
}
func (self requiredTypedRImpl) GetUint64(name string) (uint64, error) {
	v, err := self.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint64E(v)
}
func (self requiredTypedRImpl) GetByteSize(name string) (bs.ByteSize, error) {
	v, err := self.GetString(name)
	if err != nil {
		return 0, err
	}

	return bs.Parse(v)
}
