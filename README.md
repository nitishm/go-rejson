# Go-ReJSON - a golang client for ReJSON (a JSON data type for Redis)
Go-ReJSON is a [Go](https://golang.org/) client for [rejson](https://github.com/RedisLabsModules/rejson) redis module.
> ReJSON is a Redis module that implements ECMA-404 The JSON Data Interchange Standard as a native data type. It allows storing, updating and fetching JSON values from Redis keys (documents).
> Primary features:
    > Full support of the JSON standard
    > JSONPath-like syntax for selecting element inside documents
    > Documents are stored as binary data in a tree structure, allowing fast access to sub-elements
    > Typed atomic operations for all JSON values types
- Go-ReJSON is built atop the [redigo](https://github.com/gomodule/redigo) client. 
- The package is intended to be used in conjuction with the [redigo](https://github.com/gomodule/redigo), which means all features provided by the original package will be available.

## Example usage
```golang
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
```
