package main

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/janakerman/flux-event-store/internal/reconciler"
	runinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1alpha1/run"
	runreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1alpha1/run"
	tkncontroller "github.com/tektoncd/pipeline/pkg/controller"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
)

const controllerName = "flux-wait-controller"

func main() {
	sharedmain.Main(controllerName, newController)
}

func newController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	// TODO: Take server address from environment variable.
	// TODO: Correctly configure client.
	// TODO: Think about tests that test the controller end to end.
	c := &reconciler.Reconciler{
		Logger:             logrus.New(),
		EventServerAddress: "http://localhost:8080",
		Client:             http.Client{},
	}
	impl := runreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			AgentName: controllerName,
		}
	})
	c.EnqueueAfter = impl.EnqueueAfter

	runinformer.Get(ctx).Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: tkncontroller.FilterRunRef(reconciler.APIVersion, reconciler.Kind),
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	return impl
}
