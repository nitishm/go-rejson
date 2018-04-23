#!/bin/sh
set -ex
./redis-4.0.9/src/redis-server --loadmodule redis-4.0.9/rejson/src/rejson.so &