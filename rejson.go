package rejson

import "fmt"

type ReJSONHandler struct {
	clientName     string
	implementation ReJSON
}

func NewRejsonHandler() *ReJSONHandler {
	return &ReJSONHandler{clientName: "inactive"}
}

type ReJSON interface {

	JSONSet(key, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error)

	JSONGet(key, path string, opts ...JSONGetOption) (res interface{}, err error)

	JSONMGet(path string, keys ...string) (res interface{}, err error)

	JSONDel(key, path string) (res interface{}, err error)

	JSONType(key, path string) (res interface{}, err error)

	JSONNumIncrBy(key, path string, number int) (res interface{}, err error)

	JSONNumMultBy(key, path string, number int) (res interface{}, err error)

	JSONStrAppend(key, path string, jsonstring string) (res interface{}, err error)

	JSONStrLen(key, path string) (res interface{}, err error)

	JSONArrAppend(key, path string, values ...interface{}) (res interface{}, err error)

	JSONArrLen(key, path string) (res interface{}, err error)

	JSONArrPop(key, path string, index int) (res interface{}, err error)

	JSONArrIndex(key, path string, jsonValue interface{}, optionalRange ...int) (res interface{}, err error)

	JSONArrTrim(key, path string, start, end int) (res interface{}, err error)

	JSONArrInsert(key, path string, index int, values ...interface{}) (res interface{}, err error)

	JSONObjKeys(key, path string) (res interface{}, err error)

	JSONObjLen(key, path string) (res interface{}, err error)

	JSONDebug(subCmd, key, path string) (res interface{}, err error)

	JSONForget(key, path string) (res interface{}, err error)

	JSONResp(key, path string) (res interface{}, err error)

}

// JSONSet used to set a json object
// JSON.SET <key> <path> <json>
// 		 [NX | XX]
func (r *ReJSONHandler) JSONSet(key string, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONSet(key, path, obj, NX, XX)
}

// JSONGet used to get a json object
// JSON.GET <key>
//      [INDENT indentation-string]
// 		[NEWLINE line-break-string]
// 		[SPACE space-string]
// 		[NOESCAPE]
// 		[path ...]
func (r *ReJSONHandler) JSONGet(key, path string, opts ...JSONGetOption) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONGet(key, path, opts...)
}

// JSONMGet used to get path values from multiple keys
// JSON.MGET <key> [key ...] <path>
func (r *ReJSONHandler) JSONMGet(path string, keys ...string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONMGet(path, keys...)
}

// JSONDel to delete a json object
// JSON.DEL <key> <path>
func (r *ReJSONHandler) JSONDel(key string, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONDel(key, path)
}

// JSONType to get the type of key or member at path.
// JSON.TYPE <key> [path]
func (r *ReJSONHandler) JSONType(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONType(key, path)
}

// JSONNumIncrBy to increment a number by provided amount
// JSON.NUMINCRBY <key> <path> <number>
func (r *ReJSONHandler) JSONNumIncrBy(key, path string, number int) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONNumIncrBy(key, path, number)
}

// JSONNumMultBy to increment a number by provided amount
// JSON.NUMMULTBY <key> <path> <number>
func (r *ReJSONHandler) JSONNumMultBy(key, path string, number int) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONNumMultBy(key, path, number)
}

// JSONStrAppend to append a jsonstring to an existing member
// JSON.STRAPPEND <key> [path] <json-string>
func (r *ReJSONHandler) JSONStrAppend(key, path, jsonstring string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONStrAppend(key, path, jsonstring)
}

// JSONStrLen to return the length of a string member
// JSON.STRLEN <key> [path]
func (r *ReJSONHandler) JSONStrLen(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONStrLen(key, path)
}

// JSONArrAppend to append json value into array at path
// JSON.ARRAPPEND <key> <path> <json> [json ...]
func (r *ReJSONHandler) JSONArrAppend(key, path string, values ...interface{}) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrAppend(key, path, values...)
}

// JSONArrLen returns the length of the json array at path
// JSON.ARRLEN <key> [path]
func (r *ReJSONHandler) JSONArrLen(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrLen(key, path)
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
// JSON.ARRPOP <key> [path [index]]
func (r *ReJSONHandler) JSONArrPop(key, path string, index int) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrPop(key, path, index)
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
// JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
func (r *ReJSONHandler) JSONArrIndex(key, path string, jsonValue interface{}, optionalRange ...int) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrIndex(key, path, jsonValue, optionalRange...)
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
// JSON.ARRTRIM <key> <path> <start> <stop>
func (r *ReJSONHandler) JSONArrTrim(key, path string, start, end int) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrTrim(key, path, start, end)
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
// JSON.ARRINSERT <key> <path> <index> <json> [json ...]
func (r *ReJSONHandler) JSONArrInsert(key, path string, index int, values ...interface{}) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONArrInsert(key, path, index, values...)
}

// JSONObjKeys returns the keys in the object that's referenced by path
// JSON.OBJKEYS <key> [path]
func (r *ReJSONHandler) JSONObjKeys(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONObjKeys(key, path)
}

// JSONObjLen report the number of keys in the JSON Object at path in key
// JSON.OBJLEN <key> [path]
func (r *ReJSONHandler) JSONObjLen(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONObjLen(key, path)
}

// JSONDebug reports information
// JSON.DEBUG <subcommand & arguments>
//		JSON.DEBUG MEMORY <key> [path]	- report the memory usage in bytes of a value. path defaults to root if not provided.
//		JSON.DEBUG HELP					- reply with a helpful message
func (r *ReJSONHandler) JSONDebug(subCmd, key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONDebug(subCmd, key, path)
}

//JSONForget is an alias for JSONDel
//
func (r *ReJSONHandler) JSONForget(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONForget(key, path)
}

//JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//JSON.RESP <key> [path]
func (r *ReJSONHandler) JSONResp(key, path string) (res interface{}, err error) {
	if r.clientName == "inactive" {
		return nil, fmt.Errorf("no redis client is set")
	}
	return r.implementation.JSONResp(key, path)
}