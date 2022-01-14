> _Note: Currently, go-ReJSON only support redislabs/rejson with version <=1.0.8. If you are using higher versions, some commands might not work as expected_

# Go-ReJSON - a golang client for ReJSON (a JSON data type for Redis)

Go-ReJSON is a [Go](https://golang.org/) client for [ReJSON](https://github.com/RedisLabsModules/rejson) Redis Module.

[![Go Reference](https://pkg.go.dev/badge/github.com/nitishm/go-rejson.svg)](https://pkg.go.dev/github.com/nitishm/go-rejson/v4)
![test](https://github.com/nitishm/go-rejson/workflows/test/badge.svg)
![code-analysis](https://github.com/nitishm/go-rejson/workflows/code-analysis/badge.svg)
[![codecov](https://coveralls.io/repos/github/nitishm/go-rejson/badge.svg?branch=master)](https://coveralls.io/github/nitishm/go-rejson?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nitishm/go-rejson)](https://goreportcard.com/report/github.com/nitishm/go-rejson)
[![GitHub release](https://img.shields.io/github/release/nitishm/go-rejson.svg)](https://github.com/nitishm/go-rejson/releases)

> ReJSON is a Redis module that implements ECMA-404 The JSON Data Interchange Standard as a native data type. It allows storing, updating and fetching JSON values from Redis keys (documents).

Primary features of ReJSON Module:

    * Full support of the JSON standard
    * JSONPath-like syntax for selecting element inside documents
    * Documents are stored as binary data in a tree structure, allowing fast access to sub-elements
    * Typed atomic operations for all JSON values types

Each and every feature of ReJSON Module is fully incorporated in the project.

Enjoy ReJSON with the type-safe Redis client, [`Go-Redis/Redis`](https://github.com/go-redis/redis) or use the
print-like Redis-api client [`GoModule/Redigo`](https://github.com/gomodule/redigo). Go-ReJSON supports both the
clients. Use any of the above two client you want, Go-ReJSON helps you out with all its features and functionalities in
a more generic and standard way.

Support for `mediocregopher/radix` and other Redis clients is in our RoadMap. Any contributions to the support for other
clients is hearty welcome.

## Installation

    go get github.com/nitishm/go-rejson

## Example usage

```golang
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/nitishm/go-rejson/v4"
	goredis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
)

// Name - student name
type Name struct {
	First  string `json:"first,omitempty"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last,omitempty"`
}

// Student - student object
type Student struct {
	Name Name `json:"name,omitempty"`
	Rank int  `json:"rank,omitempty"`
}

func Example_JSONSet(rh *rejson.Handler) {

	student := Student{
		Name: Name{
			"Mark",
			"S",
			"Pronto",
		},
		Rank: 1,
	}
	res, err := rh.JSONSet("student", ".", student)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}

	if res.(string) == "OK" {
		fmt.Printf("Success: %s\n", res)
	} else {
		fmt.Println("Failed to Set: ")
	}

	studentJSON, err := redis.Bytes(rh.JSONGet("student", "."))
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}

	readStudent := Student{}
	err = json.Unmarshal(studentJSON, &readStudent)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Printf("Student read from redis : %#v\n", readStudent)
}

func main() {
	var addr = flag.String("Server", "localhost:6379", "Redis server address")

	rh := rejson.NewReJSONHandler()
	flag.Parse()

	// Redigo Client
	conn, err := redis.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to connect to redis-server @ %s", *addr)
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			log.Fatalf("Failed to communicate to redis-server @ %v", err)
		}
	}()
	rh.SetRedigoClient(conn)
	fmt.Println("Executing Example_JSONSET for Redigo Client")
	Example_JSONSet(rh)

	// GoRedis Client
	cli := goredis.NewClient(&goredis.Options{Addr: *addr})
	defer func() {
		if err := cli.FlushAll(context.Background()).Err(); err != nil {
			log.Fatalf("goredis - failed to flush: %v", err)
		}
		if err := cli.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
	rh.SetGoRedisClient(cli)
	fmt.Println("\nExecuting Example_JSONSET for Redigo Client")
	Example_JSONSet(rh)
}
```
