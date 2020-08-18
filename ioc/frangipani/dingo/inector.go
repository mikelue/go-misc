/*
This package provides some convenient method for utilizing dependency injection of dingo(flamingo.me).

AppContext

This interface would be used for getting instance of managed instance with panic if something gets wrong.

  appContext := AsAppContext(injector)
  instance := appContext.GetInstance(new(yourType))
  env := appContext.Environment()
*/
package dingo

import (
	"flamingo.me/dingo"
	fg "github.com/mikelue/go-misc/ioc/frangipani"
)

// Gets an "AppContext" by a "*dingo.Injector"
func AsAppContext(injector *dingo.Injector) AppContext {
	return (*appContextImpl)(injector)
}

// Main enhanced interface for "dingo"
//
// See: https://github.com/i-love-flamingo/dingo
type AppContext interface {
	// Gets instance of an object/interface
	GetInstance(interface{}) interface{}
	// Gets object of "Environment"
	//
	// See: https://pkg.go.dev/github.com/mikelue/go-misc/ioc/frangipani?tab=doc#Environment
	Environment() fg.Environment
}

type appContextImpl dingo.Injector
func (self *appContextImpl) GetInstance(of interface{}) interface{} {
	injector := (*dingo.Injector)(self)

	obj, err := injector.GetInstance(of)
	if err != nil {
		panic(err)
	}

	return obj
}
func (self *appContextImpl) Environment() fg.Environment {
	injector := (*dingo.Injector)(self)

	obj, err := injector.GetInstance(new(fg.Environment))
	if err != nil {
		panic(err)
	}

	return obj.(fg.Environment)
}
