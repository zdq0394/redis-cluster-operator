package controller

import (
	objectRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// CRD is the custom resource definition.
// Its Initialize will ensure the Custom Resource Definition present in K8S clusters.
// And it knows how to interact with k8s to List and Watch the CRD instances.
type CRD interface {
	GetListerWatcher() cache.ListerWatcher
	GetObject() objectRuntime.Object
	Initialize() error
}
