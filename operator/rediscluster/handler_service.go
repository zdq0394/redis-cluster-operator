package rediscluster

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *RedisClusterHandler) ensurePresent(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	h.Manager.EnsureRedisHeadlessService(rc, labels, ownerRefs)
	return nil
}
