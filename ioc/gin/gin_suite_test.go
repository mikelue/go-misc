package gin

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestByGinkgo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gin Suite")
}

func newContext() (*gin.Context, *httptest.ResponseRecorder) {
	resp := httptest.NewRecorder()
	sampleContext, _ := gin.CreateTestContext(resp)
	return sampleContext, resp
}
func newContextByMime(mimes ...string) *gin.Context {
	sampleContext, _ := newContext()
	sampleContext.Request = httptest.NewRequest("GET", "/car-info", nil)

	for _, mime := range mimes {
		sampleContext.Request.Header.Add("Accept", mime)
	}
	return sampleContext
}
