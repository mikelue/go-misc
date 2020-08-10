package gorm

import (
	"runtime"
	"fmt"

	"github.com/jinzhu/gorm"
)

// This exception implements "Error()/String()" interface
// to support error behaivor of GoLang
type DbException struct {
	cause error
	gormFuncName string
	callerFile string
	callerLine int
}

// Gets the file name of caller
func (self *DbException) GetCallerFile() string {
	return self.callerFile
}
// Gets the file line of caller
func (self *DbException) GetCallerLine() int {
	return self.callerLine
}
// Gets the function/feature name of GORM
func (self *DbException) GetGormFuncName() string {
	return self.gormFuncName
}
// Gets the cause of error
func (self *DbException) GetCause() error {
	return self.cause
}
// Builds the error with detailed information
func (self *DbException) DetailError() error {
	return fmt.Errorf("GORM: [%s] @ \"%s\"[%d] : %w",
		self.gormFuncName,
		self.callerFile, self.callerLine,
		self.cause,
	)
}
// As output of "DetailError()"
func (self *DbException) Error() string {
	return self.DetailError().Error()
}
// As output of "DetailError()"
func (self *DbException) String() string {
	return self.Error()
}
func (self *DbException) setFrame(frame *runtime.Frame) {
	self.callerFile = frame.File
	self.callerLine = frame.Line
}

func newDbExceptionByError(err error) *DbException {
	return &DbException{ cause: err }
}
func newDbExceptionByPanic(panicContent interface{}) *DbException {
	return newDbExceptionByError(fmt.Errorf("%v", panicContent))
}

func decoratePanic(gormFuncName string, panicContent interface{}) *DbException {
	if panicContent == nil {
		return nil
	}

	var exception *DbException
	if errorP, ok := panicContent.(error); ok {
		exception = newDbExceptionByError(errorP)
	} else {
		exception = newDbExceptionByPanic(panicContent)
	}

	exception.setFrame(getFrame(3))
	exception.gormFuncName = gormFuncName
	return exception
}
func panicError(funcName string, err error) {
	if err == nil {
		return
	}

	exception := newDbExceptionByError(err)
	exception.setFrame(getFrame(3))
	exception.gormFuncName = funcName
	panic(exception)
}
func panicIfAnyError(funcName string, sourceDb *gorm.DB) {
	var exception *DbException

	if sourceDb.RecordNotFound() {
		return
	}

	if sourceDb.Error != nil {
		exception = newDbExceptionByError(sourceDb.Error)
	} else if dbErrors := sourceDb.GetErrors(); len(dbErrors) > 0 {
		allErrors := make(gorm.Errors, 0, len(dbErrors))
		allErrors = allErrors.Add(dbErrors...)
		exception = newDbExceptionByError(sourceDb.Error)
	} else {
		return
	}

	exception.setFrame(getFrame(3))
	exception.gormFuncName = funcName

	panic(exception)
}

func panicIfAnyErrorOfAssociation(funcName string, association *gorm.Association) {
	if association.Error == nil {
		return
	}

	var exception *DbException
	exception = newDbExceptionByError(association.Error)
	exception.setFrame(getFrame(3))
	exception.gormFuncName = funcName

	panic(exception)
}

func getFrame(skip int) *runtime.Frame {
	pc := make([]uintptr, 1)
	runtime.Callers(skip + 1, pc)
	frame, _ := runtime.CallersFrames(pc).Next()

	return &frame
}
