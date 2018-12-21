package k8s

import (
	"github.com/zdq0394/redis-cluster-operator/log"
	redisclusterv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	redisclusterclientset "github.com/zdq0394/redis-cluster-operator/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// RedisCluster the RedisCluster service that knows how to interact with k8s to get them
type RedisCluster interface {
	ListRedisClusters(namespace string, opts metav1.ListOptions) (*redisclusterv1alpha1.RedisClusterList, error)
	WatchRedisClusters(namespace string, opts metav1.ListOptions) (watch.Interface, error)
}

// RedisClusterService is the RedisCluster service implementation using API calls to kubernetes.
type RedisClusterService struct {
	crdClient redisclusterclientset.Interface
	logger    log.Logger
}

// NewRedisClusterService returns a new Workspace KubeService.
func NewRedisClusterService(crdcli redisclusterclientset.Interface, logger log.Logger) *RedisClusterService {
	logger = logger.With("service", "k8s.redisfailover")
	return &RedisClusterService{
		crdClient: crdcli,
		logger:    logger,
	}
}

func (r *RedisClusterService) ListRedisClusters(namespace string, opts metav1.ListOptions) (*redisclusterv1alpha1.RedisClusterList, error) {
	return r.crdClient.Redis().RedisClusters(namespace).List(opts)
}

func (r *RedisClusterService) WatchRedisClusters(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return r.crdClient.Redis().RedisClusters(namespace).Watch(opts)
}
