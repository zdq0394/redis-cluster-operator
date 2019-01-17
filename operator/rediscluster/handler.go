package rediscluster

import (
	"context"
	"fmt"

	"github.com/zdq0394/k8soperator/pkg/util"
	manager "github.com/zdq0394/redis-cluster-operator/operator/rediscluster/handler"
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	"github.com/zdq0394/redis-cluster-operator/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultLabels = map[string]string{
	"Creator": "RedisClusterOperator",
}

const (
	RedisClusterLabelKey = "rediscluster"
)

// RedisClusterHandler handles RedisClusterCRD object to create, update or destroy a Redis Cluster
// accoring the Action and given RedisClusterCRD object and release all the resources.
type RedisClusterHandler struct {
	Labels  map[string]string
	Manager manager.RedisClusterManager
	logger  log.Logger
}

// NewRedisClusterHandler create new handler to process the watched RedisClusterCRD
func NewRedisClusterHandler(labels map[string]string, mgr manager.RedisClusterManager, logger log.Logger) *RedisClusterHandler {
	curLabels := util.MergeLabels(defaultLabels, labels)
	return &RedisClusterHandler{
		Labels:  curLabels,
		Manager: mgr,
		logger:  logger,
	}
}

// Add process the logic when a RedisClusterCRD object created/updated
// Create or update a redis cluster as RedisClusterCRD desired.
func (h *RedisClusterHandler) Add(ctx context.Context, obj runtime.Object) error {
	// Create RedisCluster Here...
	h.logger.Infoln("Create RedisCluster Here...")
	rc, ok := obj.(*redisv1alpha1.RedisCluster)
	if !ok {
		return fmt.Errorf("Cannot handle redis cluster")
	}
	h.logger.Infof("Handler Create RedisCluster:%s/%s", rc.Namespace, rc.Name)
	oRefs := h.createOwnerReferences(rc)
	instanceLabels := h.generateInstanceLabels(rc)
	labels := util.MergeLabels(h.Labels, rc.Labels, instanceLabels)
	return h.ensurePresent(rc, labels, oRefs)
}

// Delete process the logic when a RedisClusterCRD object deleted.
// Destroy the redis cluster.
func (h *RedisClusterHandler) Delete(ctx context.Context, key string) error {
	// Delete Redis Cluster
	h.logger.Infoln("Delete RedisCluster Here:", key)
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

func (h *RedisClusterHandler) generateInstanceLabels(rc *redisv1alpha1.RedisCluster) map[string]string {
	return map[string]string{
		RedisClusterLabelKey: rc.Name,
	}
}

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
