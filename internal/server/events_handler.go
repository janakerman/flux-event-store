package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/janakerman/flux-event-store/internal/storage"
)

func (s *EventServer) eventHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handlePostEvent(w, r)
	case http.MethodGet:
		s.handleGetEvents(w, r)
	default:
		s.logger.Errorf("unsupported method %s", r.Method)
		w.WriteHeader(400)
	}
}

func (s *EventServer) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	revision := params.Get("revision")
	if revision == "" {
		s.logger.Errorf("Parameter revision is required")
		w.WriteHeader(400)
	}

	events, err := s.store.EventByRevision(r.Context(), revision)
	if err != nil {
		s.logger.WithError(err).Errorf("failed to get events")
		w.WriteHeader(500)
	}

	payload, err := json.Marshal(events)
	if err != nil {
		s.logger.WithError(err).Errorf("failed to marshall events: %v", payload)
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	_, err = w.Write(payload)
	if err != nil {
		s.logger.WithError(err).Errorf("failed to write body")
		w.WriteHeader(500)
	}
}

func (s *EventServer) handlePostEvent(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Errorf("event has no revision: %s", string(body))
		w.WriteHeader(400)
	}

	err = s.store.WriteEvent(r.Context(), event)
	if err != nil {
		s.logger.WithError(err).Errorf("failed to store event")
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
}
