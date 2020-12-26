/*
Go-ReJSON is a Go client for ReJSON redis module (https://github.com/RedisLabsModules/rejson)

	ReJSON is a Redis module that aims to provide full support for ECMA-404
	The JSON Data Interchange Standard as a native data type.

	It allows storing, updating and fetching JSON values from Redis keys (documents).

	Primary features of ReJSON Module:
		* Full support of the JSON standard
		* JSONPath-like syntax for selecting element inside documents
		* Documents are stored as binary data in a tree structure, allowing fast access
		  to sub-elements
		* Typed atomic operations for all JSON values types

Go-ReJSON implements all the features of ReJSON Module, without any dependency on the client used for Redis in GoLang.

Enjoy ReJSON with the type-safe Redis client, Go-Redis/Redis (https://github.com/go-redis/redis) or use the print-like
Redis-api client GoModule/Redigo (https://github.com/gomodule/redigo/redis).

Go-ReJSON supports both the clients. Use any of the above two client you want, Go-ReJSON helps you out with all its
features and functionalities in a more generic and standard way.


Installation


To install and use ReJSON module, one must have the pre-requisites installed and setup. Run the script in :

	./install-redis-rejson.sh


Examples


Create New ReJSON Handler
	rh := rejson.NewReJSONHandler()

Set Redigo Client and use ReJSON in it
	conn, _ := redis.Dial("tcp", *addr)
	rh.SetRedigoClient(conn)

Similarly, one can set client for Go-Redis
	cli := goredis.NewClient(&goredis.Options{Addr: *addr})
	rh.SetGoRedisClient(cli)

And now, one can directly use ReJSON commands using the handler
	res, err := rh.JSONSet("str", ".", "string")

*/
package rejson
