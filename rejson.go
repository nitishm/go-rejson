package rejson

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// PopArrLast gives index of the last element for JSONArrPop
const PopArrLast = -1

// commandMux maps command name to a command function
var commandMux = map[string]func(argsIn ...interface{}) (argsOut []interface{}, err error){
	"JSON.SET":       commandJSONSet,
	"JSON.GET":       commandJSONGet,
	"JSON.DEL":       commandJSONDel,
	"JSON.MGET":      commandJSONMGet,
	"JSON.TYPE":      commandJSONType,
	"JSON.NUMINCRBY": commandJSONNumIncrBy,
	"JSON.NUMMULTBY": commandJSONNumMultBy,
	"JSON.STRAPPEND": commandJSONStrAppend,
	"JSON.STRLEN":    commandJSONStrLen,
	"JSON.ARRAPPEND": commandJSONArrAppend,
	"JSON.ARRLEN":    commandJSONArrLen,
	"JSON.ARRPOP":    commandJSONArrPop,
	"JSON.ARRINDEX":  commandJSONArrIndex,
	"JSON.ARRTRIM":   commandJSONArrTrim,
	"JSON.ARRINSERT": commandJSONArrInsert,
}

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

func commandJSONType(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONNumIncrBy(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	number := argsIn[2]
	argsOut = append(argsOut, key, path, number)
	return
}

func commandJSONNumMultBy(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	number := argsIn[2]
	argsOut = append(argsOut, key, path, number)
	return
}

func commandJSONStrAppend(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	jsonstring := argsIn[2]
	argsOut = append(argsOut, key, path, jsonstring)
	return
}

func commandJSONStrLen(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONArrAppend(argsIn ...interface{}) (argsOut []interface{}, err error) {
	keys := argsIn[0]
	path := argsIn[1]
	values := argsIn[2:]
	argsOut = append(argsOut, keys, path)
	for _, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		argsOut = append(argsOut, jsonValue)
	}
	return
}

func commandJSONArrLen(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONArrPop(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	index := argsIn[2]

	argsOut = append(argsOut, key, path)
	if index.(int) != PopArrLast {
		argsOut = append(argsOut, index)
	}
	return
}

func commandJSONArrIndex(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	jsonValue, err := json.Marshal(argsIn[2])
	if err != nil {
		return nil, err
	}
	argsOut = append(argsOut, key, path, jsonValue)

	ln := len(argsIn)
	if ln >= 4 { // start is present
		start := argsIn[3]
		argsOut = append(argsOut, start)
		if ln == 5 { // both start and end are present
			end := argsIn[4]
			argsOut = append(argsOut, end)
		}
	}
	return
}

func commandJSONArrInsert(argsIn ...interface{}) (argsOut []interface{}, err error) {
	keys := argsIn[0]
	path := argsIn[1]
	index := argsIn[2]
	values := argsIn[3:]
	argsOut = append(argsOut, keys, path, index)
	for _, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		argsOut = append(argsOut, jsonValue)
	}
	return
}

func commandJSONArrTrim(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	start := argsIn[2]
	end := argsIn[3]
	argsOut = append(argsOut, key, path, start, end)
	return
}

// CommandBuilder is used to build a command that can be used directly with redigo's conn.Do()
// This is especially useful if you do not need to conn.Do() and instead need to use the JSON.* commands in a
// MUTLI/EXEC scenario along with some other operations like GET/SET/HGET/HSET/...
func CommandBuilder(commandNameIn string, argsIn ...interface{}) (commandNameOut string, argsOut []interface{}, err error) {
	cmd, ok := commandMux[commandNameIn]
	if !ok {
		return commandNameOut, nil, fmt.Errorf("command %s not supported by ReJSON", commandNameIn)
	}

	argsOut, err = cmd(argsIn...)
	if err != nil {
		return commandNameOut, nil, fmt.Errorf("failed to execute command %s: %v", commandNameIn, err)
	}

	return commandNameIn, argsOut, nil
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
	name, args, _ := CommandBuilder("JSON.GET", key, path)
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
	name, args, _ := CommandBuilder("JSON.MGET", args...)
	return conn.Do(name, args...)
}

// JSONDel to delete a json object
// JSON.DEL <key> <path>
func JSONDel(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.DEL", key, path)
	return conn.Do(name, args...)
}

// JSONType to get the type of key or member at path.
// JSON.TYPE <key> [path]
func JSONType(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.TYPE", key, path)
	return conn.Do(name, args...)
}

// JSONNumIncrBy to increment a number by provided amount
// JSON.NUMINCRBY <key> <path> <number>
func JSONNumIncrBy(conn redis.Conn, key string, path string, number int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.NUMINCRBY", key, path, number)
	return conn.Do(name, args...)
}

// JSONNumMultBy to increment a number by provided amount
// JSON.NUMMULTBY <key> <path> <number>
func JSONNumMultBy(conn redis.Conn, key string, path string, number int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.NUMMULTBY", key, path, number)
	return conn.Do(name, args...)
}

// JSONStrAppend to append a jsonstring to an existing member
// JSON.STRAPPEND <key> [path] <json-string>
func JSONStrAppend(conn redis.Conn, key string, path string, jsonstring string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.STRAPPEND", key, path, jsonstring)
	return conn.Do(name, args...)
}

// JSONStrLen to return the length of a string member
// JSON.STRLEN <key> [path]
func JSONStrLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.STRLEN", key, path)
	return conn.Do(name, args...)
}

// JSONArrAppend to append json value into array at path
// JSON.ARRAPPEND <key> <path> <json> [json ...]
func JSONArrAppend(conn redis.Conn, key string, path string, values ...interface{}) (res interface{}, err error) {
	if len(values) == 0 {
		err = fmt.Errorf("Need atlesat one value string as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, key, path)
	args = append(args, values...)
	name, args, _ := CommandBuilder("JSON.ARRAPPEND", args...)
	return conn.Do(name, args...)
}

// JSONArrLen returns the length of the json array at path
// JSON.ARRLEN <key> [path]
func JSONArrLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRLEN", key, path)
	return conn.Do(name, args...)
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
// JSON.ARRPOP <key> [path [index]]
func JSONArrPop(conn redis.Conn, key, path string, index int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRPOP", key, path, index)
	return conn.Do(name, args...)
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
// JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
func JSONArrIndex(conn redis.Conn, key, path string, jsonValue interface{}, optionalRange ...int) (res interface{}, err error) {
	args := []interface{}{key, path, jsonValue}

	ln := len(optionalRange)
	if ln > 2 {
		return nil, fmt.Errorf("Need atmost two integeral value as an argument representing array range")
	} else if ln == 1 { // only inclusive start is present
		args = append(args, optionalRange[0])
	} else if ln == 2 { // both inclusive start and exclusive end are present
		args = append(args, optionalRange[0], optionalRange[1])
	}
	name, args, _ := CommandBuilder("JSON.ARRINDEX", args...)
	return conn.Do(name, args...)
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
// JSON.ARRTRIM <key> <path> <start> <stop>
func JSONArrTrim(conn redis.Conn, key, path string, start, end int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRTRIM", key, path, start, end)
	return conn.Do(name, args...)
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
// JSON.ARRINSERT <key> <path> <index> <json> [json ...]
func JSONArrInsert(conn redis.Conn, key, path string, index int, values ...interface{}) (res interface{}, err error) {
	if len(values) == 0 {
		err = fmt.Errorf("Need atlesat one value string as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, key, path, index)
	args = append(args, values...)
	name, args, _ := CommandBuilder("JSON.ARRINSERT", args...)
	return conn.Do(name, args...)
}
