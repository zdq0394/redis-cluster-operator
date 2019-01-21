package handler

import (
	"fmt"
)

const (
	svcHeadlessPort     = 6379
	svcHeadlessPortName = "headless-port"

	svcAccessPort     = 6379
	svcAccessPortName = "access-port"

	redisPortName = "redis"
	redisPort     = 6379

	redisClusterPortName = "cluster"
	redisClusterPort     = 16379

	graceTime = 30
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

func generateName(prefix, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}
