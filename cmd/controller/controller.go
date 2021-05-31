package main

import (
	"context"

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
	c := &reconciler.Reconciler{}
	impl := runreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			AgentName: controllerName,
		}
	})
	c.EnqueueAfter = impl.EnqueueAfter

	runinformer.Get(ctx).Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: tkncontroller.FilterRunRef("github.com/janakerman/flux-wait/v0", "FluxWait"),
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	return impl
}
