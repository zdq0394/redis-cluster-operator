package rediscluster

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zdq0394/redis-cluster-operator/operator"
	manager "github.com/zdq0394/redis-cluster-operator/operator/rediscluster/service"
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultLabels = map[string]string{
	"Creator": "RedisClusterOperator",
}

// RedisClusterHandler handles RedisClusterCRD object to create, update or destroy a Redis Cluster
// accoring the Action and given RedisClusterCRD object and release all the resources.
type RedisClusterHandler struct {
	Labels  map[string]string
	Manager manager.RedisClusterManager
}

// NewRedisClusterHandler create new handler to process the watched RedisClusterCRD
func NewRedisClusterHandler(labels map[string]string, mgr manager.RedisClusterManager) *RedisClusterHandler {
	curLabels := operator.MergeLabels(defaultLabels, labels)
	return &RedisClusterHandler{
		Labels:  curLabels,
		Manager: mgr,
	}
}

// Add process the logic when a RedisClusterCRD object created/updated
// Create or update a redis cluster as RedisClusterCRD desired.
func (h *RedisClusterHandler) Add(ctx context.Context, obj runtime.Object) error {
	// Create RedisCluster Here...
	log.Infoln("Create RedisCluster Here...")
	rc, ok := obj.(*redisv1alpha1.RedisCluster)
	if !ok {
		return fmt.Errorf("Cannot handle redis cluster")
	}
	log.Infof("Handler Create RedisCluster:%s/%s", rc.Namespace, rc.Name)
	oRefs := h.createOwnerReferences(rc)
	instanceLabels := h.generateInstanceLabels(rc)
	labels := operator.MergeLabels(h.Labels, rc.Labels, instanceLabels)
	return h.ensurePresent(rc, labels, oRefs)
}

// Delete process the logic when a RedisClusterCRD object deleted.
// Destroy the redis cluster.
func (h *RedisClusterHandler) Delete(ctx context.Context, key string) error {
	// Delete Redis Cluster
	log.Infoln("Delete RedisCluster Here:", key)
	// No need to do anything, it will be handled by the owner reference done
	// on the creation.
	return nil
}

func (h *RedisClusterHandler) createOwnerReferences(rc *redisv1alpha1.RedisCluster) []metav1.OwnerReference {
	rcvk := redisv1alpha1.VersionKind(redisv1alpha1.RCKind)
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(rc, rcvk),
	}
}

const (
	RedisClusterLabelKey = "rediscluster"
)

func (h *RedisClusterHandler) generateInstanceLabels(rc *redisv1alpha1.RedisCluster) map[string]string {
	return map[string]string{
		RedisClusterLabelKey: rc.Name,
	}
}
