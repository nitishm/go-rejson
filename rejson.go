package rejson

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func commandJSONSet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	obj := argsIn[2]
	NX := argsIn[3].(bool)
	XX := argsIn[4].(bool)
	argsOut = append(argsOut, key, path)

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	argsOut = append(argsOut, b)

	if NX && XX {
		err = fmt.Errorf("Both NX and XX cannot be true.")
		return nil, err
	} else {
		if NX {
			argsOut = append(argsOut, "NX")
		} else if XX {
			argsOut = append(argsOut, "XX")
		}
	}
	return
}

func commandJSONGet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func CommandBuilder(commandNameIn string, argsIn ...interface{}) (commandNameOut string, argsOut []interface{}, err error) {
	commandNameOut = commandNameIn
	switch commandNameIn {
	case "JSON.SET":
		argsOut, err = commandJSONSet(argsIn...)
		if err != nil {
			return "", nil, err
		}
		break
	case "JSON.GET":
		argsOut, err = commandJSONGet(argsIn...)
		if err != nil {
			return "", nil, err
		}
		break
	}
	return
}

/* JSONSet
 * JSON.SET <key> <path> <json>
 *		 [NX | XX]
**/
func JSONSet(conn redis.Conn, key string, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error) {
	name, args, err := CommandBuilder("JSON.SET", key, path, obj, NX, XX)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
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
func JSONGet(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, err := CommandBuilder("JSON.GET", key, path)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}
