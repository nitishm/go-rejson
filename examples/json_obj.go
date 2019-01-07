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

	type Object struct {
		Name      string `json:"name"`
		LastSeen  int64  `json:"lastSeen"`
		LoggedOut bool   `json:"loggedOut"`
	}

	obj := Object{"Leonard Cohen", 1478476800, true}
	res, err := rejson.JSONSet(conn, "obj", ".", obj, false, false)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}
	fmt.Println("obj:", res)

	res, err = rejson.JSONGet(conn, "obj", ".")
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	var objOut Object
	err = json.Unmarshal(res.([]byte), &objOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("got obj:", objOut)

	res, err = rejson.JSONObjLen(conn, "obj", ".")
	if err != nil {
		log.Fatalf("Failed to JSONObjLen")
		return
	}
	fmt.Println("length:", res)

	res, err = rejson.JSONObjKeys(conn, "obj", ".")
	if err != nil {
		log.Fatalf("Failed to JSONObjKeys")
		return
	}
	fmt.Println("keys:", res)

	res, err = rejson.JSONDebug(conn, rejson.DebugHelpSubcommand, "obj", ".")
	if err != nil {
		log.Fatalf("Failed to JSONDebug")
		return
	}
	fmt.Println(res)
	res, err = rejson.JSONDebug(conn, rejson.DebugMemorySubcommand, "obj", ".")
	if err != nil {
		log.Fatalf("Failed to JSONDebug")
		return
	}
	fmt.Println("Memory used by obj:", res)

	res, err = rejson.JSONGet(conn, "obj", ".",
		rejson.NewJSONGetOptionIndent("\t"), rejson.NewJSONGetOptionNewLine("\n"),
		rejson.NewJSONGetOptionSpace(" "), rejson.NewJSONGetOptionNoEscape())
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}
	err = json.Unmarshal(res.([]byte), &objOut)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}
	fmt.Println("got obj with options:", objOut)
}
