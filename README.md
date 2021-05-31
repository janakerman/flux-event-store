# flux-event-store

A proof of concept of a Flux event API combined with a Tekton Custom Task controller. The custom run waits for a given
Flux revision to reconcile before marking the task as complete.

## Usage

Install Tekton pipelines:

```
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.24.1/release.yaml
```

Install the POC components:
```
kubectl apply -f <(helm template ./charts/flux-event-store)
```

Create a custom Task:
```
apiVersion: tekton.dev/v1alpha1
kind: Run
metadata:
  name: run-my-custom-task-xjns6
spec:
  ref:
    apiVersion: github.com/janakerman/flux-event-store/v0
    kind: FluxDeploymentWait
  params:
  - name: revision
    value: b
```

