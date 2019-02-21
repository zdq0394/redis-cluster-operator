package k8s

import (
	k8sservice "github.com/zdq0394/k8soperator/pkg/k8s"
	"github.com/zdq0394/k8soperator/pkg/log"
	redisclusterclientset "github.com/zdq0394/redis-cluster-operator/pkg/client/clientset/versioned"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// Service is the K8s service entrypoint.
type Services interface {
	k8sservice.CRD
	k8sservice.ConfigMap
	k8sservice.Pod
	k8sservice.PodDisruptionBudget
	RedisCluster
	k8sservice.Service
	k8sservice.RBAC
	k8sservice.Deployment
	k8sservice.StatefulSet
}

type services struct {
	k8sservice.CRD
	k8sservice.ConfigMap
	k8sservice.Pod
	k8sservice.PodDisruptionBudget
	RedisCluster
	k8sservice.Service
	k8sservice.RBAC
	k8sservice.Deployment
	k8sservice.StatefulSet
}

// New returns a new Kubernetes service.
func New(kubecli kubernetes.Interface, crdcli redisclusterclientset.Interface, apiextcli apiextensionscli.Interface, logger log.Logger) Services {
	return &services{
		CRD:                 k8sservice.NewCRDService(apiextcli, logger),
		ConfigMap:           k8sservice.NewConfigMapService(kubecli, logger),
		Pod:                 k8sservice.NewPodService(kubecli, logger),
		PodDisruptionBudget: k8sservice.NewPodDisruptionBudgetService(kubecli, logger),
		RedisCluster:        NewRedisClusterService(crdcli, logger),
		Service:             k8sservice.NewServiceService(kubecli, logger),
		RBAC:                k8sservice.NewRBACService(kubecli, logger),
		Deployment:          k8sservice.NewDeploymentService(kubecli, logger),
		StatefulSet:         k8sservice.NewStatefulSetService(kubecli, logger),
	}
}
