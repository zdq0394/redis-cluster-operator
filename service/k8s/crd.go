package k8s

import (
	"github.com/spotahome/kooper/client/crd"
)

// CRDConf is alias of crd.Conf
type CRDConf = crd.Conf

// CRD is the CRD service that knows how to interact with k8s to manage them.
type CRD interface {
	// CreateWorkspaceCRD will create the custom resource and wait to be ready.
	EnsureCRD(conf CRDConf) error
}

// CRDService is the CRD service implementation using API calls to kubernetes.
type CRDService struct {
	crdCli crd.interface
	logger log.Logger
}

// NewCRDService returns a new CRD KubeService.
func NewCRDService(aeClient apiextensionscli.Interface, logger log.Logger) *CRDService {
	logger = logger.With("service", "k8s.crd")
	crdCli := crd.NewClient(aeClient, logger)

	return &CRDService{
		crdCli: crdCli,
		logger: logger,
	}
}

