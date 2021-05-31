# flux-event-store

A proof of concept of a Flux event API combined with a Tekton Custom Task controller. The custom run waits for a given
Flux revision to reconcile before marking the task as complete.

It comprises of two components:
1. An event API that can be configured as a Flux notification provider to provide a persistent event store.
2. A Tekton Custom Task controller which uses the event API to wait for Flux Kustomizations to reconcile successfully before
marking the task as successful. 

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
    value: a
```

### A failed Flux release

1. Create a custom Task waiting for revision `a`

    ```
    apiVersion: tekton.dev/v1alpha1
    kind: Run
    metadata:
      name: run-my-custom-task-1
    spec:
      ref:
        apiVersion: github.com/janakerman/flux-event-store/v0
        kind: FluxDeploymentWait
      params:
      - name: revision
        value: a
    ```

2. Exec onto a pod and fake some Flux notifications (replace the IP of the pod `flux-event-store`):

    ```
    curl -X POST -vvv 10.244.0.39:8080/events \
        -d '{"metadata":{"revision":"a"},"timestamp":"0001-01-01T00:01:00Z","reason":"BuildFailed"}'
    ```
   
3. The `Run` fails.

    ```
    $ kubectl get runs
    NAME                       SUCCEEDED   REASON                 STARTTIME   COMPLETIONTIME
    run-my-custom-task-1   False       ReconciliationFailed   7m25s       5m15s
    ```
   
### A successful Flux release

1. Create a custom Task waiting for revision `b`

    ```
cat <<EOF | kubectl apply -f -
apiVersion: tekton.dev/v1alpha1
kind: Run
metadata:
  name: run-my-custom-task-2
spec:
  ref:
    apiVersion: github.com/janakerman/flux-event-store/v0
    kind: FluxDeploymentWait
  params:
  - name: revision
    value: b
EOF
    ```

2. Exec onto a pod and fake some Flux notifications (replace the IP of the pod `flux-event-store`):

    ```
    curl -X POST -vvv 10.244.0.39:8080/events \
        -d '{"metadata":{"revision":"b"},"timestamp":"0001-01-01T00:01:00Z","reason":"ReconciliationSucceeded"}'
    ```
   
3. The `Run` fails.

    ```
    $ kubectl get runs
      NAME                       SUCCEEDED   REASON                   STARTTIME   COMPLETIONTIME
      run-my-custom-task-2       True        ReconciliationSuceeded   13m         33s
    ```
   
## TODO:

- [ ] Return events from more reconciler paths? These seem to show up in the Events when describing the Run which is helpful
for debugging.
- [ ] Tests for reconciler.
- [ ] End-to-end tests - manual testing taking far too long.
- [ ] Persistent storage for api.
- [ ] Decide whether to split into multiple repos or leave as is.
- [ ] Sort out build
     - [ ] Add correct tag into Helm Chart (copy podinfo)
     - [ ] Build images on release/tags