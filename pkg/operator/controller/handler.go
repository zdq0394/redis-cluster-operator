package controller

import (
	"context"

	objectRuntime "k8s.io/apimachinery/pkg/runtime"
)

// Handler knows how to handle the received resources from a kubernetes cluster.
type Handler interface {
	Add(context.Context, objectRuntime.Object) error
	Delete(context.Context, string) error
}
