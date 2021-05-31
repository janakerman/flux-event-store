package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janakerman/flux-event-store/internal/storage"

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
			req := httptest.NewRequest("GET", "/anything", bytes.NewBuffer(nil))
			res := httptest.NewRecorder()
			server.handleHealth(res, req)

			Expect(res.Code).To(Equal(200))
		})
	})

	Context("event endpoint", func() {
		payload := func(event storage.Event) *bytes.Buffer {
			json, err := json.Marshal(event)
			Expect(err).To(BeNil())
			return bytes.NewBuffer(json)
		}

		var (
			store   storage.EventStore
			handler http.HandlerFunc
		)

		BeforeEach(func() {
			store = storage.NewInMemory()
			handler = server.eventHandler(store)
		})

		It("stores event", func() {
			event := storage.Event{
				Metadata: storage.EventMetaData{
					Revision: "branch/sha",
				},
			}

			req := httptest.NewRequest("GET", "/anything", payload(event))
			res := httptest.NewRecorder()
			handler(res, req)

			Expect(res.Code).To(Equal(200))
			events, err := store.EventByRevision(context.Background(), event.Metadata.Revision)
			Expect(err).To(BeNil())
			Expect(events).To(HaveLen(1))
			Expect(events[0]).To(Equal(event))
		})
		It("returns 400 if event has no revision", func() {
			event := storage.Event{
				Metadata: storage.EventMetaData{
					Revision: "",
				},
			}

			req := httptest.NewRequest("GET", "/anything", payload(event))
			res := httptest.NewRecorder()
			handler(res, req)

			Expect(res.Code).To(Equal(400))
		})
	})
})
