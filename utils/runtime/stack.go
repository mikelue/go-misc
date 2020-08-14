/*
This package contains utilities of manipulating objects on package of "runtime".

CallerUtils

This space provides convenient methods to retrieve information from "runtime.Caller()"
*/
package runtime

import (
	"runtime"
	"path/filepath"
)

// Method space for runtime.Caller.
var CallerUtils ICallerUtils = 0

type ICallerUtils int
// Gets the dictionary of current source file.
func (ICallerUtils) GetDirOfSource() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
