package rediscluster

import (
	"context"
	"fmt"

	"github.com/zdq0394/k8soperator/pkg/log"
	"github.com/zdq0394/k8soperator/pkg/util"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster/cluster"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster/sentinel"
	v1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultLabels = map[string]string{
	"Creator": "RedisClusterOperator",
}

const (
	RedisClusterLabelKey = "rediscluster"
	RedisClusterCluster  = "cluster"
	RedisClusterSentinel = "sentinel"
)

// Handler handles RedisClusterCRD object to create, update or destroy a Redis Cluster
// accoring the Action and given RedisClusterCRD object and release all the resources.
type Handler struct {
	Labels          map[string]string
	ClusterManager  cluster.RedisClusterManager
	SentinelManager sentinel.RedisSentinelManager
	logger          log.Logger
}

// NewRedisClusterHandler create new handler to process the watched RedisClusterCRD
func NewRedisClusterHandler(labels map[string]string, clusterMgr cluster.RedisClusterManager, sentinelMgr sentinel.RedisSentinelManager, logger log.Logger) *Handler {
	curLabels := util.MergeLabels(defaultLabels, labels)
	return &Handler{
		Labels:          curLabels,
		ClusterManager:  clusterMgr,
		SentinelManager: sentinelMgr,
		logger:          logger,
	}
}

// Add process the logic when a RedisClusterCRD object created/updated
// Create or update a redis cluster as RedisClusterCRD desired.
func (h *Handler) Add(ctx context.Context, obj runtime.Object) error {
	// Create RedisCluster Here...
	h.logger.Infoln("Create RedisCluster Here...")
	rc, ok := obj.(*v1alpha1.RedisCluster)
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
func (h *Handler) Delete(ctx context.Context, key string) error {
	// Delete Redis Cluster
	h.logger.Infoln("Delete RedisCluster Here:", key)
	// No need to do anything, it will be handled by the owner reference done
	// on the creation.
	return nil
}

func (h *Handler) createOwnerReferences(rc *v1alpha1.RedisCluster) []metav1.OwnerReference {
	rcvk := v1alpha1.VersionKind(v1alpha1.RCKind)
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(rc, rcvk),
	}
}

func (h *Handler) generateInstanceLabels(rc *v1alpha1.RedisCluster) map[string]string {
	return map[string]string{
		RedisClusterLabelKey: rc.Name,
	}
}

func (h *Handler) ensurePresent(rc *v1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	if rc.Spec.Mode == RedisClusterCluster {
		return h.ensureClusterPresent(rc, labels, ownerRefs)
	} else if rc.Spec.Mode == RedisClusterSentinel {
		return h.ensureSentinelPresent(rc, labels, ownerRefs)
	}
	return fmt.Errorf("invalid redis cluster mode:%s. valid modes are:[%s,%s]", rc.Spec.Mode, RedisClusterCluster, RedisClusterSentinel)
}

func (h *Handler) ensureClusterPresent(rc *v1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	// Create Redis ConfigMap
	if err := h.ClusterManager.EnsureRedisConfigMap(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Headless service for statefulset
	if err := h.ClusterManager.EnsureRedisHeadlessService(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Statefulset
	if err := h.ClusterManager.EnsureRedisStatefulset(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Redis Access Service
	if err := h.ClusterManager.EnsureRedisAcessService(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Wait Redis Statefulset Pods is Running
	if err := h.ClusterManager.WaitRedisStatefulsetPodsRunning(rc, labels, ownerRefs); err != nil {
		return err
	}
	// Create Boot pod to create redis cluster
	if err := h.ClusterManager.EnsureRedisClusterBootPod(rc, labels, ownerRefs); err != nil {
		return err
	}
	return nil
}

func (h *Handler) ensureSentinelPresent(rc *v1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) error {
	if err := h.SentinelManager.EnsureRedisConfigMap(rc, labels, ownerRefs); err != nil {
		return err
	}
	if err := h.SentinelManager.EnsureRedisHeadlessService(rc, labels, ownerRefs); err != nil {
		return err
	}
	if err := h.SentinelManager.EnsureRedisStatefulset(rc, labels, ownerRefs); err != nil {
		return err
	}
	return nil
}
