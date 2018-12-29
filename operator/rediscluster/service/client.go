package service

import (
	"fmt"
	"time"

	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/service/k8s"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisClusterManager manage redis cluster instance in kubernetes cluster
type RedisClusterManager interface {
	EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, lables map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	EnsureRedisAcessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
	WaitRedisStatefulsetPodsRunning(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error
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

func (s *redisKubeClusterManager) EnsureRedisConfigMap(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	configMap := generateRedisConfigMap(rc, labels, ownerRefs)
	return s.K8SService.CreateOrUpdateConfigMap(rc.Namespace, configMap)
}

func (s *redisKubeClusterManager) EnsureRedisStatefulset(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	ss := generateRedisStatefulset(rc, labels, ownerRefs)
	return s.K8SService.CreateOrUpdateStatefulSet(rc.Namespace, ss)
}

func (s *redisKubeClusterManager) WaitRedisStatefulsetPodsRunning(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	name := generateName(statefulsetNamePrefix, rc.Name)
	namespace := rc.Namespace
	var err error
	var ss *appsv1beta2.StatefulSet
	<-time.After(time.Duration(6) * time.Minute)
	timeOut := time.After(time.Duration(12) * time.Minute)
	t := time.NewTicker(time.Duration(1) * time.Minute)
	for {
		select {
		case <-t.C:
			ss, err = s.K8SService.GetStatefulSet(namespace, name)
			if err != nil {
				fmt.Println("Statefulset Replicas:", ss.Status.Replicas)
				fmt.Println("Statefulset ReadyReplicas:", ss.Status.ReadyReplicas)
				fmt.Println("Statefulset CurrentReplicas:", ss.Status.CurrentReplicas)
				if ss.Status.Replicas == ss.Status.ReadyReplicas && ss.Status.Replicas == ss.Status.CurrentReplicas {
					return nil
				}
			}
		case <-timeOut:
			return fmt.Errorf("Timeout waiting for Statefulset to be Running")
		}
	}
}

func (s *redisKubeClusterManager) EnsureRedisHeadlessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	svc := generateRedisHeadlessService(rc, labels, ownerRefs)
	return s.K8SService.CreateIfNotExistsService(rc.Namespace, svc)
}

func (s *redisKubeClusterManager) EnsureRedisAcessService(rc *redisv1alpha1.RedisCluster, labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	svc := generateRedisAccessService(rc, labels, ownerRefs)
	return s.K8SService.CreateIfNotExistsService(rc.Namespace, svc)
}
