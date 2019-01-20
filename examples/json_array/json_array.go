package main

/*
import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

func main() {
	var addr = flag.String("Server", "localhost:6379", "Redis server address")

	flag.Parse()

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

	res, err = rejson.JSONArrIndex(conn, "arr", ".", "one")
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("Index of \"one\":", res)

	res, err = rejson.JSONArrIndex(conn, "arr", ".", "three", 3, 10)
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("Out of range:", res)

	res, err = rejson.JSONArrIndex(conn, "arr", ".", "ten")
	if err != nil {
		log.Fatalf("Failed to JSONArrIndex %v", err)
		return
	}
	fmt.Println("\"ten\" not found:", res)

	res, err = rejson.JSONArrTrim(conn, "arr", ".", 1, 2)
	if err != nil {
		log.Fatalf("Failed to JSONArrTrim %v", err)
		return
	}
	fmt.Println("no. of elements left:", res)

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
	fmt.Println("arr after trimming to [1,2]:", ArrOut)

	res, err = rejson.JSONArrInsert(conn, "arr", ".", 0, "one")
	if err != nil {
		log.Fatalf("Failed to JSONArrInsert %v", err)
		return
	}
	fmt.Println("no. of elements:", res)

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
	fmt.Println("arr after inserting \"one\":", ArrOut)
}
*/