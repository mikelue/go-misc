/*
This package contains various utilities to ease development of GoLang.

EnvExecutor

This object can execute your routine with customized environment variables.

The "os.Environ()" would be reverted to original status after the execution.

  yourEnvExec := NewEnvExecutor(map[string]string{ "HOME": "./test-home" })
  // The environment is un-changed after the execution
  yourEnvExec.Execute(your_func)
*/
package utils

import (
	"os"
)

// Constructs an new executor with expected environment.
//
// This object cannot be used in multi-thread situation.
func NewEnvExecutor(changedEnvVars map[string]string) *EnvExecutor {
	return &EnvExecutor{
		changedEnvVars,
		map[string]string{},
	}
}
//
type EnvExecutor struct {
	neededVars map[string]string
	oldVars map[string]string
}
// Executes a routine in modified environment.
func (self *EnvExecutor) Run(callback func()) {
	defer func() {
		/**
		 * Reverts the old value of environment variable
		 */
		for k, v := range self.oldVars {
			os.Setenv(k, v)
		}
		// :~)

		/**
		 * Un-set the environment variables
		 */
		for k, _ := range self.neededVars {
			if _, existedInOld := self.oldVars[k]; !existedInOld {
				os.Unsetenv(k)
			}
		}
		// :~)

		self.oldVars = map[string]string{}
	}()

	/**
	 * Set environment variable and backup old value
	 */
	for k, v := range self.neededVars {
		if oldValue, ok := os.LookupEnv(k); ok {
			self.oldVars[k] = oldValue
		}

		os.Setenv(k, v)
	}
	// :~)

	callback()
}
