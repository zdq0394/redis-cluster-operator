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
)

var (
	terminationGracePeriodSeconds int64 = 20
	storageClassNamePx                  = "px-hdd-ha3"
)
