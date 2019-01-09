FROM centos:7.4
COPY make/release/redisops /opt/redis/operator

WORKDIR /opt/redis

ENTRYPOINT ["operator", "cluster"]