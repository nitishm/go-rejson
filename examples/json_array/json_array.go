package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/nitishm/go-rejson/v4/rjs"

	"github.com/gomodule/redigo/redis"
	goredis "github.com/redis/go-redis/v9"

	"github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()

func Example_JSONArray(rh *rejson.Handler) {
	ArrIn := []string{"one", "two", "three", "four", "five"}
	res, err := rh.JSONSet("arr", ".", ArrIn)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}
	fmt.Println("arr:", res)

	res, err = rh.JSONGet("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	var ArrOut []string
	err = json.Unmarshal(res.([]byte), &ArrOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("arr before pop:", ArrOut)

	res, err = rh.JSONArrLen("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONArrLen")
		return
	}
	fmt.Println("Length:", res)

	res, err = rh.JSONArrPop("arr", ".", rjs.PopArrLast)
	if err != nil {
		log.Fatalf("Failed to JSONArrLen")
		return
	}
	var ele string
	err = json.Unmarshal(res.([]byte), &ele)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("Deleted element:", ele)

	res, err = rh.JSONGet("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	err = json.Unmarshal(res.([]byte), &ArrOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("arr after pop:", ArrOut)

	res, err = rh.JSONArrLen("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONArrLen")
		return
	}
	fmt.Println("Length:", res)

	res, err = rh.JSONArrIndex("arr", ".", "one")
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("Index of \"one\":", res)

	res, err = rh.JSONArrIndex("arr", ".", "three", 3, 10)
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("Out of range:", res)

	res, err = rh.JSONArrIndex("arr", ".", "ten")
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("\"ten\" not found:", res)

	res, err = rh.JSONArrTrim("arr", ".", 1, 2)
	if err != nil {
		log.Fatalf("Failed to JSONArrTrim %v", err)
		return
	}
	fmt.Println("no. of elements left:", res)

	res, err = rh.JSONGet("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	err = json.Unmarshal(res.([]byte), &ArrOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("arr after trimming to [1,2]:", ArrOut)

	res, err = rh.JSONArrInsert("arr", ".", 0, "one")
	if err != nil {
		log.Fatalf("Failed to JSONArrInsert %v", err)
		return
	}
	fmt.Println("no. of elements:", res)

	res, err = rh.JSONGet("arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	err = json.Unmarshal(res.([]byte), &ArrOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("arr after inserting \"one\":", ArrOut)
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
	Example_JSONArray(rh)

	// GoRedis Client
	cli := goredis.NewClient(&goredis.Options{Addr: *addr})
	defer func() {
		if err := cli.FlushAll(ctx).Err(); err != nil {
			log.Fatalf("goredis - failed to flush: %v", err)
		}
		if err := cli.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
	rh.SetGoRedisClientWithContext(context.Background(), cli)
	fmt.Println("\nExecuting Example_JSONSET for Redigo Client")
	Example_JSONArray(rh)
}
