package storage

import (
	"context"
	"sort"
	"time"
)

// Example Kustomization event.
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

// Failed Kustomization event:
// {
//  "involvedObject": {
//    "kind": "Kustomization",
//    "namespace": "flux-system",
//    "name": "flux-system",
//    "uid": "7e114434-c5ab-4ba7-a151-877a82b4312b",
//    "apiVersion": "kustomize.toolkit.fluxcd.io/v1beta1",
//    "resourceVersion": "69472"
//  },
//  "severity": "error",
//  "timestamp": "2021-05-31T13:58:25Z",
//  "message": "kustomize create failed: failed to decode Kubernetes YAML from /tmp/flux-system673152982/clusters/my-cluster/podinfo-helm-release.yaml: error converting YAML to JSON: yaml: line 13: could not find expected ':'",
//  "reason": "BuildFailed",
//  "metadata": {
//    "revision": "main/4857e9edb52f87859b4af8b4edbc05b594cd7a69"
//  },
//  "reportingController": "kustomize-controller",
//  "reportingInstance": "kustomize-controller-6b7578445c-jr2ss"
//}
//
//{
//  "involvedObject": {
//    "kind": "Kustomization",
//    "namespace": "flux-system",
//    "name": "flux-system",
//    "uid": "7e114434-c5ab-4ba7-a151-877a82b4312b",
//    "apiVersion": "kustomize.toolkit.fluxcd.io/v1beta1",
//    "resourceVersion": "77336"
//  },
//  "severity": "error",
//  "timestamp": "2021-05-31T14:32:23Z",
//  "message": "validation failed: error: error validating \"7e114434-c5ab-4ba7-a151-877a82b4312b.yaml\": error validating data: ValidationError(HelmRelease.spec.chart.spec): missing required field \"chart\" in io.fluxcd.toolkit.helm.v2beta1.HelmRelease.spec.chart.spec; if you choose to ignore these errors, turn validation off with --validate=false\n",
//  "reason": "ValidationFailed",
//  "metadata": {
//    "revision": "main/5d888586c0e85c8f98862d2122a44c052ab394e0"
//  },
//  "reportingController": "kustomize-controller",
//  "reportingInstance": "kustomize-controller-6b7578445c-jr2ss"
//}

type (
	Events []Event
	Event  struct {
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

func (e Event) HasSuceeded() bool {
	return e.Reason == "ReconciliationSucceeded"
}

func (e Event) HasFailed() bool {
	return e.Reason == "BuildFailed" || e.Reason == "ValidationFailed"
}

func (e Events) MostRecent() *Event {
	if len(e) == 0 {
		return nil
	}
	sort.Sort(EventsByTimestamp(e))

	return &e[len(e)-1]
}

type EventsByTimestamp Events

func (e EventsByTimestamp) Len() int {
	return len(e)
}

func (e EventsByTimestamp) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e EventsByTimestamp) Less(i, j int) bool {
	return e[i].Timestamp.Before(e[j].Timestamp)
}
