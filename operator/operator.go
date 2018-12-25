package operator

import (
	"github.com/zdq0394/redis-cluster-operator/log"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster"
	k8sclient "github.com/zdq0394/redis-cluster-operator/pkg/k8s"
	"github.com/zdq0394/redis-cluster-operator/pkg/operator"
	"github.com/zdq0394/redis-cluster-operator/pkg/operator/controller"
	k8service "github.com/zdq0394/redis-cluster-operator/service/k8s"
)

// Start the Operator
func Start(development bool, kubeconfig string) error {
	kubeClient, redisClient, aeClient, _ := k8sclient.CreateKubernetesClients(development, kubeconfig)
	logger := log.Base()
	kubeService := k8service.New(kubeClient, redisClient, aeClient, logger)
	crd := rediscluster.NewRedisClusterCRD(kubeService)
	redisClusterHandler := rediscluster.NewRedisClusterHandler()
	ctrl := controller.NewSimpleController(crd, redisClusterHandler)
	optor := operator.NewSimpleOperator(crd, ctrl)
	stopC := make(chan struct{}, 0)
	optor.Run(stopC)
	return nil
}
