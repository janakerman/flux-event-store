package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/janakerman/flux-event-store/internal/storage"

	"github.com/sirupsen/logrus"
)

type EventServer struct {
	port   string
	logger *logrus.Logger
}

func NewEventServer(port string, logger *logrus.Logger) *EventServer {
	return &EventServer{
		port:   port,
		logger: logger,
	}
}

// ListenAndServe starts the HTTP server on the specified port
func (s *EventServer) ListenAndServe(stopCh <-chan struct{}) {
	mux := http.DefaultServeMux
	mux.Handle("/health", http.HandlerFunc(s.handleHealth))
	mux.Handle("/notification", http.HandlerFunc(s.eventHandler(storage.NewInMemory())))
	srv := &http.Server{
		Addr:    s.port,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error(err, "Event server crashed")
			os.Exit(1)
		}
	}()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error(err, "Event server graceful shutdown failed")
	} else {
		s.logger.Info("Event server stopped")
	}
}

func (s *EventServer) eventHandler(store storage.EventStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.WithError(err).Errorf("failed to read request body")
			w.WriteHeader(500)
		}

		event := storage.Event{}

		err = json.Unmarshal(body, &event)
		if err != nil {
			s.logger.WithError(err).Errorf("failed to unmarshall request body: %s", string(body))
			w.WriteHeader(400)
		}

		if event.Metadata.Revision == "" {
			s.logger.WithError(err).Errorf("event has no revision: %s", string(body))
			w.WriteHeader(400)
		}

		err = store.WriteEvent(r.Context(), event)
		if err != nil {
			s.logger.WithError(err).Errorf("failed to store event")
			w.WriteHeader(500)
		}

		w.WriteHeader(200)
	}
}

func (s *EventServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}
