package storage

import (
	"context"
	"time"
)

// Example event.
//
// {
//  "involvedObject": {
//    "kind": "Kustomization",
//    "namespace": "flux-system",
//    "name": "flux-system",
//    "uid": "7e114434-c5ab-4ba7-a151-877a82b4312b",
//    "apiVersion": "kustomize.toolkit.fluxcd.io/v1beta1",
//    "resourceVersion": "37155"
//  },
//  "severity": "info",
//  "timestamp": "2021-05-31T11:49:10Z",
//  "message": "Update completed",
//  "reason": "ReconciliationSucceeded",
//  "metadata": {
//    "commit_status": "update",
//    "revision": "main/594a4f770b226c3acf29736fdf2c0e835b7e037f"
//  },
//  "reportingController": "kustomize-controller",
//  "reportingInstance": "kustomize-controller-6b7578445c-jr2ss"
//}
type (
	Event struct {
		Timestamp       time.Time
		Message, Reason string
		Metadata        EventMetaData
	}
	EventMetaData struct {
		CommitStatus string `json:"commit_status"`
		Revision     string
	}
)

type EventStore interface {
	WriteEvent(ctx context.Context, event Event) error
	EventByRevision(ctx context.Context, revision string) ([]Event, error)
}
