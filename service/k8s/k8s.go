package k8s

import (
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	"github.com/zdq0394/redis-cluster-operator/log"
	redisclusterclientset "github.com/zdq0394/redis-cluster-operator/pkg/client/clientset/versioned"
)

// Service is the K8s service entrypoint.
type Services interface {
	CRD
	ConfigMap
	Pod
	PodDisruptionBudget
	RedisCluster
	Service
	RBAC
	Deployment
	StatefulSet
}

type services struct {
	CRD
	ConfigMap
	Pod
	PodDisruptionBudget
	RedisCluster
	Service
	RBAC
	Deployment
	StatefulSet
}

// New returns a new Kubernetes service.
func New(kubecli kubernetes.Interface, crdcli redisclusterclientset.Interface, apiextcli apiextensionscli.Interface, logger log.Logger) Services {
	return &services{
		CRD:                 NewCRDService(apiextcli, logger),
		ConfigMap:           NewConfigMapService(kubecli, logger),
		Pod:                 NewPodService(kubecli, logger),
		PodDisruptionBudget: NewPodDisruptionBudgetService(kubecli, logger),
		RedisCluster:        NewRedisClusterService(crdcli, logger),
		Service:             NewServiceService(kubecli, logger),
		RBAC:                NewRBACService(kubecli, logger),
		Deployment:          NewDeploymentService(kubecli, logger),
		StatefulSet:         NewStatefulSetService(kubecli, logger),
	}
}
