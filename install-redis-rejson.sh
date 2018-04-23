#!/bin/sh
set -ex
wget http://download.redis.io/releases/redis-4.0.9.tar.gz
tar -xzvf redis-4.0.9.tar.gz 
cd redis-4.0.9 && make 

git clone https://github.com/RedisLabsModules/rejson.git
cd rejson && make