package main

import (
	"fmt"
	rejson "go-rejson"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Name struct {
	First  string `json:"first,omitempty"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last,omitempty"`
}

type Student struct {
	Name Name `json:"name,omitempty"`
	Rank int  `json:"rank,omitempty"`
}

func main() {
	conn, err := redis.Dial("tcp", "6390")
	if err != nil {
		log.Fatal("Failed to connect to port 6390")
	}

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
		return
	}

	fmt.Printf("Success if OK - %s\n", res)
}
