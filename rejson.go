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
		err = fmt.Errorf("Both NX and XX cannot be true")
		return nil, err
	}

	if NX {
		argsOut = append(argsOut, "NX")
	} else if XX {
		argsOut = append(argsOut, "XX")
	}
	return
}

func commandJSONGet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONDel(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONMGet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	keys := argsIn[0 : len(argsIn)-1]
	path := argsIn[len(argsIn)-1]
	argsOut = append(argsOut, keys...)
	argsOut = append(argsOut, path)
	return
}

// CommandBuilder is used to build a command that can be used directly with redigo's conn.Do()
// This is especially useful if you do not need to conn.Do() and instead need to use the JSON.* commands in a
// MUTLI/EXEC scenario along with some other operations like GET/SET/HGET/HSET/...
func CommandBuilder(commandNameIn string, argsIn ...interface{}) (commandNameOut string, argsOut []interface{}, err error) {
	commandNameOut = commandNameIn
	switch commandNameIn {
	case "JSON.SET":
		argsOut, err = commandJSONSet(argsIn...)
		if err != nil {
			return "", nil, err
		}
	case "JSON.GET":
		argsOut, err = commandJSONGet(argsIn...)
		if err != nil {
			return "", nil, err
		}
	case "JSON.DEL":
		argsOut, err = commandJSONDel(argsIn...)
		if err != nil {
			return "", nil, err
		}
	case "JSON.MGET":
		argsOut, err = commandJSONMGet(argsIn...)
		if err != nil {
			return "", nil, err
		}
	default:
		err = fmt.Errorf("Command %s not supported by ReJSON", commandNameIn)
		return "", nil, err
	}
	return
}

// JSONSet used to set a json object
// JSON.SET <key> <path> <json>
// 		 [NX | XX]
func JSONSet(conn redis.Conn, key string, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error) {
	name, args, err := CommandBuilder("JSON.SET", key, path, obj, NX, XX)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}

// JSONGet used to get a json object
// JSON.GET <key>
//      [INDENT indentation-string]
// 	[NEWLINE line-break-string]
// 		[SPACE space-string]
// 	[NOESCAPE]
// 	[path ...]
func JSONGet(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, err := CommandBuilder("JSON.GET", key, path)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}

// JSONMGet used to get path values from multiple keys
// JSON.MGET <key> [key ...] <path>
func JSONMGet(conn redis.Conn, path string, keys ...string) (res interface{}, err error) {
	if len(keys) == 0 {
		err = fmt.Errorf("Need atlesat one key as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	args = append(args, path)
	name, args, err := CommandBuilder("JSON.MGET", args...)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}

// JSONDel to delete a json object
// JSON.DEL <key> <path>
func JSONDel(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, err := CommandBuilder("JSON.DEL", key, path)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}
