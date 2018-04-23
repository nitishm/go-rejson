package rejson

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

/* JSONSet
 * JSON.SET <key> <path> <json>
 *		 [NX | XX]
**/
func JSONSet(conn redis.Conn, key string, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error) {
	var command []interface{}
	command = append(command, key, path)

	b, err := json.Marshal(obj)
	if err != nil {
		return
	}
	command = append(command, b)

	if NX && XX {
		err = fmt.Errorf("Both NX and XX cannot be true.")
		return
	} else {
		if NX {
			command = append(command, "NX")
		} else if XX {
			command = append(command, "XX")
		}
	}
	return conn.Do("JSON.SET", command...)
}

/**
 * JSON.GET <key>
 *      [INDENT indentation-string]
 *		[NEWLINE line-break-string]
 * 		[SPACE space-string]
 *		[NOESCAPE]
 *		[path ...]
 * The args dont make much sense in terms of the client.
**/
func JSONGet(conn redis.Conn, key string, path string, args ...interface{}) (res interface{}, err error) {
	var command []interface{}
	command = append(command, key, path)

	return conn.Do("JSON.GET", command...)
}
