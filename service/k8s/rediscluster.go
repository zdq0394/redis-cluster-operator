package k8s
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	redisclusterv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	redisclusterclientset "github.com/zdq0394/redis-cluster-operator/pkg/client/clientset/versioned"
)

// RedisCluster the RedisCluster service that knows how to interact with k8s to get them
type RedisCluster interface {
	ListRedisClusters(namespace string, opts metav1.ListOptions)(* redisclusterv1alpha1.RedisClusterList, error)
	WatchRedisClusters(namespace string, opts metav1.ListOptions)(* watch.interface, error)
}

// RedisClusterService is the RedisCluster service implementation using API calls to kubernetes.
type RedisClusterService struct {
	crdClient redisfailoverclientset.Interface
	logger    log.Logger
}

// NewRedisClusterService returns a new Workspace KubeService.
func NewRedisClusterService(crdcli redisclusterclientset.Interface, logger log.Logger) *RedisFailoverService {
	logger = logger.With("service", "k8s.redisfailover")
	return &RedisFailoverService{
		crdClient: crdcli,
		logger:    logger,
	}
}

func (r *RedisClusterService)ListRedisClusters(namespace string, opts metav1.ListOptions)(* redisclusterv1alpha1.RedisClusterList, error){
	return r.crdClient.Redis().RedisClusters(namespace).List(opts)
}

func (r *RedisClusterService)WatchRedisClusters(namespace string, opts metav1.ListOptions)(* watch.interface, error){
	return r.crdClient.Redis().RedisClusters(namespace).Watch(opts)
}
