package main

import (
	"encoding/json"
	"flag"
	rejson "go-rejson"
	"log"

	"github.com/gomodule/redigo/redis"
)

var addr = flag.String("Server", "localhost:6379", "Redis server address")

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

func main() {
	flag.Parse()

	conn, err := redis.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to connect to redis-server @ %s", *addr)
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	student := Student{
		Name: Name{
			"Mark",
			"S",
			"Pronto",
		},
		Rank: 1,
	}
	res, err := rejson.JSONSet(conn, "student", ".", student, false, false)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}

	log.Printf("Success if - %s\n", res)

	studentJSON, err := redis.Bytes(rejson.JSONGet(conn, "student", ""))
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

	log.Printf("Student read from redis : %#v\n", readStudent)
}
