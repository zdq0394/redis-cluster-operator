# ! /usr/bin/sh

ifneq ($(shell uname), Darwin)
    EXTLDFLAGS = -extldflags "-static" $(null)
else
    EXTLDFLAGS =
endif

BUILD_NUMBER=$(shell git rev-parse --short HEAD)

BUILD_DATE=$(shelldate +%FT%T%z)

build:
	mkdir -p make/release
	go build -o make/release/redisops github.com/zdq0394/redis-cluster-operator/cmd/operator
	chmod 775 make/release/redisops

run:
	make/release/redisops cluster --develop --kubeconfig=/root/.kube/config
