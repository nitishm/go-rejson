#!/bin/sh
set -ex
./redis-4.0.9/src/redis-server --loadmodule rejson/src/rejson.so