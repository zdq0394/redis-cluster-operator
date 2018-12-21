package crd

import (
	"time"

	kubeversion "k8s.io/kubernetes/pkg/util/version"
)

const (
	checkCRDInterval = 2 * time.Second
	crdReadyTimeout  = 3 * time.Minute
)

var (
	clusterMinVersion = kubeversion.MustParseGeneric("v1.7.0")
)
