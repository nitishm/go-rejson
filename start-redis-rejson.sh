#!/bin/sh
set -ex
./redis-4.0.9/src/redis-server --loadmodule github.com/RedisLabsModules/rejson/src/rejson.so