package service

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/service/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisClusterManager manage redis cluster instance in kubernetes cluster
type RedisClusterManager interface {
	EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, lables map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisAcessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
}

type redisKubeClusterManager struct {
	K8SService k8s.Services
}

// NewRedisClusterManager new redis cluster manager.
func NewRedisClusterManager(s k8s.Services) RedisClusterManager {
	return &redisKubeClusterManager{
		K8SService: s,
	}
}

func (s *redisKubeClusterManager) EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, lables map[string]string, ownerRefs []metav1.OwnerReference) error {
	return nil
}

func (s *redisKubeClusterManager) EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	return nil
}

func (s *redisKubeClusterManager) EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	svc := generateRedisHeadlessService(rc, labels, ownerRefs)
	return s.K8SService.CreateIfNotExistsService(rc.Namespace, svc)
}

func (s *redisKubeClusterManager) EnsureRedisAcessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	svc := generateRedisAccessService(rc, labels, ownerRefs)
	return s.K8SService.CreateIfNotExistsService(rc.Namespace, svc)
}
