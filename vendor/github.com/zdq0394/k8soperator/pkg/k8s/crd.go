package k8s

import (
	"github.com/zdq0394/k8soperator/pkg/log"
	"github.com/zdq0394/k8soperator/pkg/operator/crd"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// CRD is the CRD service that knows how to interact with k8s to manage them.
type CRD interface {
	// CreateWorkspaceCRD will create the custom resource and wait to be ready.
	EnsureCRD(conf crd.Conf) error
}

// CRDService is the CRD service implementation using API calls to kubernetes.
type CRDService struct {
	crdCli crd.Client
	logger log.Logger
}

// NewCRDService returns a new CRD KubeService.
func NewCRDService(aeClient apiextensionscli.Interface, logger log.Logger) CRD {
	logger = logger.With("service", "k8s.crd")
	crdCli := crd.NewClient(aeClient)
	return &CRDService{
		crdCli: crdCli,
		logger: logger,
	}
}

// EnsureCRD ...
func (s *CRDService) EnsureCRD(conf crd.Conf) error {
	return s.crdCli.EnsurePresent(conf)
}
