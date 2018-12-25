package rediscluster

import (
	"github.com/zdq0394/redis-cluster-operator/log"
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
	crd := NewRedisClusterCRD(kubeService)
	handler := NewRedisClusterHandler()

	cfg := &controller.Config{
		Name: "Redis Cluster Controller",
	}
	ctrl := controller.NewSimpleController(cfg, crd, handler)
	optor := operator.NewSimpleOperator(crd, ctrl)
	stopC := make(chan struct{}, 0)
	optor.Run(stopC)
	return nil
}
