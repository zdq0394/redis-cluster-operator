package rediscluster

import (
	"github.com/zdq0394/k8soperator/pkg/operator"
	"github.com/zdq0394/k8soperator/pkg/operator/controller"
	manager "github.com/zdq0394/redis-cluster-operator/operator/rediscluster/cluster"
	k8sclient "github.com/zdq0394/redis-cluster-operator/pkg/k8s"
	"github.com/zdq0394/redis-cluster-operator/pkg/log"
	k8service "github.com/zdq0394/redis-cluster-operator/service/k8s"
)

type Config struct {
	Development       bool
	Kubeconfig        string
	BootImg           string
	ClusterDomain     string
	ConcurrentWorkers int
}

// Start the Operator
func Start(conf *Config) error {
	kubeClient, redisClient, aeClient, _ := k8sclient.CreateKubernetesClients(conf.Development, conf.Kubeconfig)
	logger := log.Base()
	kubeService := k8service.New(kubeClient, redisClient, aeClient, logger)
	crd := NewRedisClusterCRD(kubeService)

	mgr := manager.NewRedisClusterManager(kubeService, conf.BootImg, conf.ClusterDomain)
	handler := NewRedisClusterHandler(nil, mgr, logger)

	controllerCfg := &controller.Config{
		Name:              "Redis Cluster Controller",
		ConcurrentWorkers: conf.ConcurrentWorkers,
	}
	ctrl := controller.NewSimpleController(controllerCfg, crd, handler)
	optor := operator.NewSimpleOperator(crd, ctrl)
	stopC := make(chan struct{}, 0)
	optor.Run(stopC)
	return nil
}
