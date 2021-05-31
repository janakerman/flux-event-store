package storage

import (
	"context"
)

var _ EventStore = &InMemory{}

type InMemory struct {
	byRevision map[string][]Event
}

func NewInMemory() *InMemory {
	return &InMemory{byRevision: map[string][]Event{}}
}

func (m *InMemory) WriteEvent(_ context.Context, event Event) error {
	m.byRevision[event.Metadata.Revision] = append(m.byRevision[event.Metadata.Revision], event)
	return nil
}

func (m *InMemory) EventByRevision(ctx context.Context, revision string) ([]Event, error) {
	return m.byRevision[revision], nil
}
