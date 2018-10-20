package rejson

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

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
	"JSON.ARRINDEX":  commandJSONArrIndex,
	"JSON.ARRPOP":    commandJSONArrPop,
	"JSON.ARRTRIM":   commandJSONArrTrim,
	"JSON.OBJKEYS":   commandJSONObjKeys,
	"JSON.OBJLEN":    commandJSONObjLen,
	"JSON.DEBUG":     commandJSONDebug,
	"JSON.FORGET":    commandJSONDel,
	"JSON.RESP":      commandJSONResp,
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
	options := argsIn[2:]
	argsOut = append(argsOut, key)
	argsOut = append(argsOut, options...)
	argsOut = append(argsOut, path)
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

func commandJSONArrIndex(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	scalar := argsIn[2]
	start := argsIn[3]
	stop := argsIn[4]
	argsOut = append(argsOut, key, path, scalar, start, stop)
	return
}

func commandJSONArrPop(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	index := argsIn[2]
	argsOut = append(argsOut, key, path, index)
	return
}


func commandJSONArrTrim(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	start := argsIn[2]
	stop := argsIn[3]
	argsOut = append(argsOut, key, path, start, stop)
	return
}

func commandJSONObjKeys(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONObjLen(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
	return
}

func commandJSONDebug(argsIn ...interface{})(argsOut []interface{}, err error) {
	subCommand := argsIn[0]
	arguments := argsIn[1:]
	argsOut = append(argsOut,subCommand)
	argsOut = append(argsOut,arguments...)
	return
}

func commandJSONResp(argsIn ...interface{})(argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key, path)
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
func JSONGet(conn redis.Conn, key string, path string, options ...string) (res interface{}, err error) {
	args := make([]interface{}, 0)
	args=append(args, key)
	args=append(args, path)
	for _, option := range options {
		args = append(args, option)
	}
	name, args, _ := CommandBuilder("JSON.GET", args...)
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

// // JSONArrLen returns the length of the json array at path
// // JSON.ARRLEN <key> [path]
func JSONArrLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRLEN", key, path)
	return conn.Do(name, args...)
}

//JSON.ARRINDEX return the position of the scalar value in the array, or -1 if unfound.
//JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
//The optional inclusive start (default 0) and exclusive stop (default 0, meaning that the last element is included) specify a slice of the array to search.
func JSONArrIndex(conn redis.Conn, key string, path string, scalar interface{}, start int, stop int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRINDEX", key, path, scalar, start, stop)
	return conn.Do(name,args...)
}

//JSON.ARRPOP Remove and return element from the index in the array.
//JSON.ARRPOP <key> [path [index]]
//index is the position in the array to start popping from (defaults to -1, meaning the last element).
func JSONArrPop(conn redis.Conn, key string, path string, index int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRPOP", key, path,index)
	return conn.Do(name,args...)
}

//JSON.ARRTRIM Trim an array so that it contains only the specified inclusive range of elements.
//JSON.ARRTRIM <key> <path> <start> <stop>
func JSONArrTrim(conn redis.Conn, key string, path string, start int, stop int) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.ARRTRIM", key, path, start, stop)
	return conn.Do(name,args...)
}

//JSON.OBJKEYS Return the keys in the object that's referenced by path.
//JSON.OBJKEYS <key> [path]
func JSONObjKeys(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.OBJKEYS", key, path)
	return conn.Do(name,args...)
}

//JSON.OBJLEN Report the number of keys in the JSON Object at path in key.
//JSON.OBJLEN  <key> [path]
func JSONObjLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.OBJLEN", key, path)
	return conn.Do(name,args...)
}

//JSON.DEBUG Supported subcommands are:
//MEMORY <key> [path] - report the memory usage in bytes of a value. path defaults to root if not provided.
// returns an integer, specifically the size in bytes of the value
//HELP - reply with a helpful message.returns an array, specifically with the help message
//JSON.DEBUG <subcommand & arguments>
func JSONDebug(conn redis.Conn,subCommand string,arguments ...string) (res interface{}, err error) {
	args := make([]interface{}, 0)
	args=append(args, subCommand)
	for _, option := range arguments {
		args = append(args, option)
	}
	name, args, _ := CommandBuilder("JSON.DEBUG",args...)
	return conn.Do(name,args...)
}

//JSON.FORGET An alias for JSON.DEL.
func JSONForget(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.FORGET", key, path)
	return conn.Do(name, args...)
}

//JSON.RESP Return the JSON in key in Redis Serialization Protocol (RESP).
//JSON.RESP <key> [path]
func JSONResp(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.RESP", key, path)
	return conn.Do(name,args...)
}