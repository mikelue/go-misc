/*
You can register "ErrorHandler" to resolve errors comes from "MvcHandler"

See
  "MvcConfig.RegisterErrorHandlers(ErrorHandler)"
*/
package gin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Process and resolves the generated error by MvcHandler
type ErrorHandler interface {
	// Checks whether or not the error bould be resolved
	CanHandle(*gin.Context, error) bool
	// Handles the error
	HandleError(*gin.Context, error) error
}

type errorController []ErrorHandler
func (self errorController) handle(context *gin.Context, err error) {
	for i, errorHandler := range self {
		if errorHandler.CanHandle(context, err) {
			if err := errorHandler.HandleError(context, err); err != nil {
				fmt.Println("Severe Error:")
				fmt.Printf("Handle error[Index %d] has failed: %v\n", i, err)
			}
			break
		}
	}
}
