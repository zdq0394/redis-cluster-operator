package sentinel

import (
	"fmt"

	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/service/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisSentinelManager manage redis sentinel mode cluster in kubernetes cluster
type RedisSentinelManager interface {
	EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, lables map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
}

type redisKubeSentinelManager struct {
	ClusterDomain string
	K8SService    k8s.Services
}

// NewRedisSentinelManager return new Sentinel Redis manager.
func NewRedisSentinelManager(s k8s.Services, clusterDomain string) RedisSentinelManager {
	return &redisKubeSentinelManager{
		ClusterDomain: clusterDomain,
		K8SService:    s,
	}
}

func (s *redisKubeSentinelManager) EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	fmt.Println("Create Redis Sentinel Configmap")
	return nil
}

func (s *redisKubeSentinelManager) EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, lables map[string]string, ownerRefs []metav1.OwnerReference) error {
	fmt.Println("Create Redis Sentinel Statefulset")
	return nil
}

func (s *redisKubeSentinelManager) EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	fmt.Println("Create Redis Sentinel Headless Service")
	return nil
}
