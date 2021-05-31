package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/janakerman/flux-event-store/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kreconciler "knative.dev/pkg/reconciler"
)

const (
	APIVersion    = "github.com/janakerman/flux-event-store/v0"
	Kind          = "FluxDeploymentWait"
	ParamRevision = "revision"
)

type EnqueueFunc func(interface{}, time.Duration)

type params struct {
	revision string
}

type Reconciler struct {
	EventServerAddress string
	EnqueueAfter       EnqueueFunc
	Logger             *logrus.Logger
	Client             http.Client
}

// ReconcileKind implements Interface.ReconcileKind.
func (c *Reconciler) ReconcileKind(ctx context.Context, r *v1alpha1.Run) kreconciler.Event {
	if r.Spec.Ref == nil ||
		r.Spec.Ref.APIVersion != APIVersion || r.Spec.Ref.Kind != Kind {
		// This is not a Run we should have been notified about; do nothing.
		return nil
	}

	if r.Spec.Ref.Name != "" {
		c.Logger.Errorf("Found unexpected ref name: %s", r.Spec.Ref.Name)
		r.Status.MarkRunFailed("UnexpectedName", "Found unexpected ref name: %s", r.Spec.Ref.Name)
		return nil
	}

	// Ignore completed waits.
	if r.IsDone() {
		c.Logger.Infof("Run is finished, done reconciling\n")
		return nil
	}

	params, err := getParams(r)
	if err != nil {
		c.Logger.WithError(err).Errorf("Unexpected parameters")
		r.Status.MarkRunFailed("UnexpectedParams", err.Error())
		return nil
	}

	// The sole purpose of this task is to wait - consider it started as of the first reconciliation loop.
	if r.Status.StartTime == nil {
		now := metav1.NewTime(time.Now())
		r.Status.StartTime = &now
	}

	u, err := url.Parse(fmt.Sprintf("%s/events", c.EventServerAddress))
	if err != nil {
		c.Logger.WithError(err).Errorf("Found unexpected ref name: %s", r.Spec.Ref.Name)
		return nil
	}
	u.Query().Add("revision", params.revision)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		c.Logger.WithError(err).Errorf("failed building request")
		return nil
	}

	res, err := c.Client.Do(req)
	if err != nil {
		c.Logger.WithError(err).Errorf("events request failed")
		return nil
	}
	defer func() {
		_ = res.Body.Close()
	}()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		c.Logger.WithError(err).Errorf("read response body")
		return nil
	}

	var events storage.Events
	err = json.Unmarshal(b, &events)
	if err != nil {
		c.Logger.WithError(err).Errorf("unmarshall events response")
		return nil
	}

	r.Status.MarkRunRunning("Waiting", "Waiting for Kustomisation to complete reconciliation")

	if len(events) == 0 {
		c.Logger.Info("no events received")
		return nil
	}

	mostRecent := events.MostRecent()

	if mostRecent.HasSuceeded() {
		c.Logger.Info("reconciliation suceeded")
		r.Status.MarkRunSucceeded("ReconciliationSuceeded", "The revision has been reconciled.")
	} else if mostRecent.HasFailed() {
		c.Logger.Info("reconciliation failed")
		r.Status.MarkRunFailed("ReconciliationFailed",
			"The revision failed to reconcile. Reason: '%s'. Message: '%s'",
			mostRecent.Reason,
			mostRecent.Message)
	} else {
		c.EnqueueAfter(r, 5*time.Second)
	}

	return kreconciler.NewEvent(corev1.EventTypeNormal, "RunReconciled", "Run reconciled: \"%s/%s\"", r.Namespace, r.Name)
}

func getParams(r *v1alpha1.Run) (*params, error) {
	revision, err := getParam(r, ParamRevision)
	if err != nil {
		return nil, err
	}

	if len(r.Spec.Params) != 1 {
		var found []string
		for _, p := range r.Spec.Params {
			if p.Name == ParamRevision {
				continue
			}
			found = append(found, p.Name)
		}
		if len(found) > 0 {
			return nil, fmt.Errorf("found unexpected params: %v", found)
		}
	}
	return &params{
		revision: revision,
	}, nil
}

func getParam(r *v1alpha1.Run, param string) (string, error) {
	expr := r.Spec.GetParam(param)
	if expr == nil || expr.Value.StringVal == "" {
		return "", fmt.Errorf("%s param is required", param)
	}
	return expr.Value.StringVal, nil
}
