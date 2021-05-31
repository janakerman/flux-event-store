package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	var (
		server EventServer
		store  storage.EventStore
	)

	BeforeEach(func() {
		store = storage.NewInMemory()
		server = EventServer{
			listenAddr: "8080",
			logger:     logrus.New(),
			store:      store,
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

	Context("event ingest endpoint", func() {
		marshallEvent := func(event storage.Event) *bytes.Buffer {
			json, err := json.Marshal(event)
			Expect(err).To(BeNil())
			return bytes.NewBuffer(json)
		}
		unmarshallEvents := func(buffer *bytes.Buffer) []storage.Event {
			b, err := io.ReadAll(buffer)
			Expect(err).To(BeNil())

			var e []storage.Event
			err = json.Unmarshal(b, &e)
			Expect(err).To(BeNil())
			return e
		}

		Context("posting an event", func() {
			It("stores event", func() {
				event := storage.Event{
					Metadata: storage.EventMetaData{
						Revision: "branch/sha",
					},
				}

				req := httptest.NewRequest("POST", "/anything", marshallEvent(event))
				res := httptest.NewRecorder()
				server.eventHandler(res, req)

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

				req := httptest.NewRequest("POST", "/anything", marshallEvent(event))
				res := httptest.NewRecorder()
				server.eventHandler(res, req)

				Expect(res.Code).To(Equal(400))
			})
		})

		Context("reading an event", func() {
			It("returns empty events", func() {
				req := httptest.NewRequest("GET", "/anything?revision=rev", nil)
				res := httptest.NewRecorder()
				server.eventHandler(res, req)

				Expect(res.Code).To(Equal(200))

				events := unmarshallEvents(res.Body)
				Expect(events).To(HaveLen(0))
			})

			It("returns the correct event", func() {
				event := storage.Event{
					Metadata: storage.EventMetaData{
						Revision: "rev",
					},
				}
				err := store.WriteEvent(context.Background(), event)
				Expect(err).To(BeNil())

				req := httptest.NewRequest("GET", "/anything?revision=rev", nil)
				res := httptest.NewRecorder()
				server.eventHandler(res, req)

				Expect(res.Code).To(Equal(200))

				events := unmarshallEvents(res.Body)
				Expect(events).To(HaveLen(1))
				Expect(events[0]).To(Equal(event))
			})

			It("returns 400 when revision param is omitted", func() {
				req := httptest.NewRequest("GET", "/anything", nil)
				res := httptest.NewRecorder()
				server.eventHandler(res, req)

				Expect(res.Code).To(Equal(400))
			})
		})
	})
})
