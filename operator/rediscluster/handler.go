package rediscluster

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RedisClusterHandler handles RedisClusterCRD object to create, update or destroy a Redis Cluster
// accoring the Action and given RedisClusterCRD object and release all the resources.
type RedisClusterHandler struct {
}

// NewRedisClusterHandler create new handler to process the watched RedisClusterCRD
func NewRedisClusterHandler() *RedisClusterHandler {
	return &RedisClusterHandler{}
}

// Add process the logic when a RedisClusterCRD object created/updated
// Create or update a redis cluster as RedisClusterCRD desired.
func (h *RedisClusterHandler) Add(ctx context.Context, obj runtime.Object) error {
	// Create RedisCluster Here...
	fmt.Println("Create RedisCluster Here...")
	return nil
}

// Delete process the logic when a RedisClusterCRD object deleted.
// Destroy the redis cluster.
func (h *RedisClusterHandler) Delete(ctx context.Context, key string) error {
	// Delete Redis Cluster
	fmt.Println("Delete RedisCluster Here:", key)
	return nil
}

func (h *RedisClusterHandler) createOwnerReferences(rc *redisv1alpha1.RedisCluster) []metav1.OwnerReference {
	rcvk := redisv1alpha1.VersionKind(redisv1alpha1.RCKind)
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(rc, rcvk),
	}
}
