#!/usr/bin/env bash

# dockerfile use
/usr/bin/redis-server /etc/redis.conf &

cd /usr/local/redis_orm_workbench/

./redis_orm_workbench

