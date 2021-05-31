package server

import (
	"bytes"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server suite")
}

var _ = Describe("Event server", func() {
	var server EventServer

	BeforeEach(func() {
		server = EventServer{
			port:   "8080",
			logger: logrus.New(),
		}
	})

	Context("health endpoint", func() {
		It("returns 200", func() {
			req := httptest.NewRequest("GET", "/health", bytes.NewBuffer(nil))
			res := httptest.NewRecorder()
			server.handleHealth(res, req)

			Expect(res.Code).To(Equal(200))
		})
	})
})
