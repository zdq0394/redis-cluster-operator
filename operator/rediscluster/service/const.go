package service

const (
	redisHeadlessPort     = 6379
	redisHeadlessPortName = "redis-port"

	redisAccessPort     = 6379
	redisAccessPortName = "redis-port"
)

const (
	svcHeadlessNamePrefix = "rch"
	svcAccessNamePrefix   = "rca"
	statefulsetNamePrefix = "rcs"
	configMapNamePrefix   = "rcfg"
	bootPodNamePrefix     = "rcboot"
)

var (
	terminationGracePeriodSeconds int64 = 20
)

const (
	redisClusterBootImage         = "zdq0394/redis-cluster-boot:1.2"
	redisClusterBootClusterDomain = "cluster.local"
)
