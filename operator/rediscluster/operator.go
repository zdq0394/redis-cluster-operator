package rediscluster

import (
	"github.com/zdq0394/k8soperator/pkg/log"
	"github.com/zdq0394/k8soperator/pkg/operator"
	"github.com/zdq0394/k8soperator/pkg/operator/controller"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster/cluster"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster/sentinel"
	k8sclient "github.com/zdq0394/redis-cluster-operator/pkg/k8s"
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

	clusterMgr := cluster.NewRedisClusterManager(kubeService, conf.BootImg, conf.ClusterDomain)
	sentinelMgr := sentinel.NewRedisSentinelManager(kubeService, conf.ClusterDomain)
	handler := NewRedisClusterHandler(nil, clusterMgr, sentinelMgr, logger)

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
