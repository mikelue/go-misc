package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP service", func() {
	const listen = ":12380"

	Context("New HTTP service", func() {
		It("Start/Stop a HTTP service", func() {
			sampleServer := &http.Server {
				Addr: listen,
				Handler: http.HandlerFunc(sampleHttpHandler),
			}

			testedService := HttpServiceBuilder.New(
				sampleServer,
				func(server *http.Server) error {
					return server.ListenAndServe()
				},
			)

			var startError error

			/**
			 * Starts HTTP service
			 */
			By("Starting HTTP Service")
			go func() {
				startError = testedService.Start(context.TODO())
			}()
			<-time.After(2 * time.Second)
			// :~)

			/**
			 * Try to make a real request to HTTP server
			 */
			resp, err := resty.New().R().
				Get(fmt.Sprintf("http://localhost%s", listen))
			Expect(err).To(Succeed())
			Expect(resp.StatusCode()).To(BeEquivalentTo(http.StatusUseProxy))
			Expect(resp.Header().Get("h1")).To(BeEquivalentTo("v1"))
			// :~)

			/**
			 * Stopping HTTP Service
			 */
			By("Stopping HTTP Service")
			Expect(testedService.Stop(context.TODO())).
				To(Succeed())
			// :~)

			/**
			 * Asserts the starting result has no error(exclude ErrServerClosed)
			 */
			Expect(startError).Should(Succeed())
			// :~)
		})
	})
})

func sampleHttpHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("h1", "v1")
	writer.WriteHeader(http.StatusUseProxy)
}
