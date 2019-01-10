FROM centos:7.4.1708
COPY make/release/redisops /opt/redis/operator

WORKDIR /opt/redis

ENTRYPOINT ["operator", "cluster"]
CMD ["--bootimg=zdq0394/redis-cluster-boot:1.2", "--clusterdomain=cluster.local", "--concurrentworkers=3"]
