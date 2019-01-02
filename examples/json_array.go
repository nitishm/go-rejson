package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"log"
)

func main() {
	var addr = flag.String("Server", "localhost:6379", "Redis server address")

	flag.Parse()

	conn, err := redis.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to connect to redis-server @ %s", *addr)
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	ArrIn := []string{"one", "two", "three", "four", "five"}
	res, err := rejson.JSONSet(conn, "arr", ".", ArrIn, false, false)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}
	fmt.Println("arr:", res)

	res, err = rejson.JSONGet(conn, "arr", ".")
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

	res, err = rejson.JSONArrLen(conn, "arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONArrLen")
		return
	}
	fmt.Println("Length:", res)

	res, err = rejson.JSONArrPop(conn, "arr", ".", rejson.PopArrLast)
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

	res, err = rejson.JSONGet(conn, "arr", ".")
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

	res, err = rejson.JSONArrLen(conn, "arr", ".")
	if err != nil {
		log.Fatalf("Failed to JSONArrLen")
		return
	}
	fmt.Println("Length:", res)

}
