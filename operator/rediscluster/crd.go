package rediscluster

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/pkg/operator/crd"
	"github.com/zdq0394/redis-cluster-operator/service/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RedisClusterCRD struct {
	service k8s.Services
}

// NewRedisClusterCRD create an instance of RedisClusterCRD
func NewRedisClusterCRD(service k8s.Services) *RedisClusterCRD {
	return &RedisClusterCRD{
		service: service,
	}
}

// Initialize ensure RedisClusterCRD created in k8s.
func (s *RedisClusterCRD) Initialize() error {
	crdConf := crd.Conf{
		Kind:       redisv1alpha1.RCKind,
		NamePlural: redisv1alpha1.RCNamePlural,
		Group:      redisv1alpha1.SchemeGroupVersion.Group,
		Version:    redisv1alpha1.SchemeGroupVersion.Version,
		Scope:      redisv1alpha1.RCScope,
		Categories: []string{"all"},
	}
	return s.service.EnsureCRD(crdConf)
}

// GetListerWatcher get the listwatcher of RedisClusterCFD.
func (s *RedisClusterCRD) GetListerWatcher() cache.ListerWatcher {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return s.service.ListRedisClusters("", options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return s.service.WatchRedisClusters("", options)
		},
	}
}

// GetObject get the RedisCluster Object
func (s *RedisClusterCRD) GetObject() runtime.Object {
	return &redisv1alpha1.RedisCluster{}
}
