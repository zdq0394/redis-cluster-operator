package rediscluster

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *RedisClusterHandler) ensurePresent(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	// Create Redis ConfigMap
	if err := h.Manager.EnsureRedisConfigMap(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Headless service for statefulset
	if err := h.Manager.EnsureRedisHeadlessService(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Statefulset
	if err := h.Manager.EnsureRedisStatefulset(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Access Service
	if err := h.Manager.EnsureRedisAcessService(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Wait Redis Statefulset Pods is Running
	if err := h.Manager.WaitRedisStatefulsetPodsRunning(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Boot pod to create redis cluster
	if err := h.Manager.EnsureRedisClusterBootPod(rc, labels, ownerRefs); err != nil {
		return err
	}
	return nil
}
