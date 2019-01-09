FROM docker-registry.saicstack.com/library/ubuntu:16.04
COPY make/release/redisops /opt/redis/operator

WORKDIR /opt/redis

ENTRYPOINT ["operator", "cluster"]