module github.com/janakerman/flux-event-store

go 1.16

require (
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/onsi/ginkgo v1.16.3
	github.com/onsi/gomega v1.13.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tektoncd/pipeline v0.24.1
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/sys v0.0.0-20210531080801-fdfd190a6549 // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b // indirect
	knative.dev/pkg v0.0.0-20210528203030-47dfdcfaedfd
	sigs.k8s.io/structured-merge-diff/v4 v4.1.1 // indirect
)

replace (
    k8s.io/client-go => k8s.io/client-go v0.20.2
    k8s.io/api => k8s.io/api v0.20.2
)
