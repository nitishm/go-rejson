package rejson

import (
	"github.com/nitishm/go-rejson/v4/rjs"
)

type Handler struct {
	clientName     string
	implementation ReJSON
}

func NewReJSONHandler() *Handler {
	return &Handler{clientName: rjs.ClientInactive}
}

// ReJSON provides an interface for various Go Redis Clients to implement ReJSON commands
type ReJSON interface {
	JSONSet(key, path string, obj interface{}, opts ...rjs.SetOption) (res interface{}, err error)

	JSONGet(key, path string, opts ...rjs.GetOption) (res interface{}, err error)

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

	JSONDebug(subCmd rjs.DebugSubCommand, key, path string) (res interface{}, err error)

	JSONForget(key, path string) (res interface{}, err error)

	JSONResp(key, path string) (res interface{}, err error)
}

// JSONSet used to set a json object
//
// ReJSON syntax:
// 	JSON.SET <key> <path> <json>
// 			 [NX | XX]
//
func (r *Handler) JSONSet(key string, path string, obj interface{}, opts ...rjs.SetOption) (
	res interface{}, err error) {

	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONSet(key, path, obj, opts...)
}

// JSONGet used to get a json object
//
// ReJSON syntax:
// 	JSON.GET <key>
//			[INDENT indentation-string]
//			[NEWLINE line-break-string]
//			[SPACE space-string]
//			[NOESCAPE]
//			[path ...]
//
func (r *Handler) JSONGet(key, path string, opts ...rjs.GetOption) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONGet(key, path, opts...)
}

// JSONMGet used to get path values from multiple keys
//
// ReJSON syntax:
// 	JSON.MGET <key> [key ...] <path>
//
func (r *Handler) JSONMGet(path string, keys ...string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONMGet(path, keys...)
}

// JSONDel to delete a json object
//
// ReJSON syntax:
// 	JSON.DEL <key> <path>
//
func (r *Handler) JSONDel(key string, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONDel(key, path)
}

// JSONType to get the type of key or member at path.
//
// ReJSON syntax:
// 	JSON.TYPE <key> [path]
//
func (r *Handler) JSONType(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONType(key, path)
}

// JSONNumIncrBy to increment a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMINCRBY <key> <path> <number>
//
func (r *Handler) JSONNumIncrBy(key, path string, number int) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONNumIncrBy(key, path, number)
}

// JSONNumMultBy to increment a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMMULTBY <key> <path> <number>
//
func (r *Handler) JSONNumMultBy(key, path string, number int) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONNumMultBy(key, path, number)
}

// JSONStrAppend to append a jsonstring to an existing member
//
// ReJSON syntax:
// 	JSON.STRAPPEND <key> [path] <json-string>
//
func (r *Handler) JSONStrAppend(key, path, jsonstring string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONStrAppend(key, path, jsonstring)
}

// JSONStrLen to return the length of a string member
//
// ReJSON syntax:
// 	JSON.STRLEN <key> [path]
//
func (r *Handler) JSONStrLen(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONStrLen(key, path)
}

// JSONArrAppend to append json value into array at path
//
// ReJSON syntax:
// 	JSON.ARRAPPEND <key> <path> <json> [json ...]
//
func (r *Handler) JSONArrAppend(key, path string, values ...interface{}) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrAppend(key, path, values...)
}

// JSONArrLen returns the length of the json array at path
//
// ReJSON syntax:
// 	JSON.ARRLEN <key> [path]
//
func (r *Handler) JSONArrLen(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrLen(key, path)
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
//
// ReJSON syntax:
// 	JSON.ARRPOP <key> [path [index]]
//
func (r *Handler) JSONArrPop(key, path string, index int) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrPop(key, path, index)
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
//
// ReJSON syntax:
// 	JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
//
func (r *Handler) JSONArrIndex(key, path string, jsonValue interface{}, optionalRange ...int) (
	res interface{}, err error) {

	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrIndex(key, path, jsonValue, optionalRange...)
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
//
// ReJSON syntax:
// 	JSON.ARRTRIM <key> <path> <start> <stop>
//
func (r *Handler) JSONArrTrim(key, path string, start, end int) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrTrim(key, path, start, end)
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
//
// ReJSON syntax:
// 	JSON.ARRINSERT <key> <path> <index> <json> [json ...]
//
func (r *Handler) JSONArrInsert(key, path string, index int, values ...interface{}) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONArrInsert(key, path, index, values...)
}

// JSONObjKeys returns the keys in the object that's referenced by path
//
// ReJSON syntax:
// 	JSON.OBJKEYS <key> [path]
//
func (r *Handler) JSONObjKeys(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONObjKeys(key, path)
}

// JSONObjLen report the number of keys in the JSON Object at path in key
//
// ReJSON syntax:
// 	JSON.OBJLEN <key> [path]
//
func (r *Handler) JSONObjLen(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONObjLen(key, path)
}

// JSONDebug reports information
//
// ReJSON syntax:
// 	JSON.DEBUG <subcommand & arguments>
//		JSON.DEBUG MEMORY <key> [path]	- report the memory usage in bytes of a value. path defaults to root if not provided.
//		JSON.DEBUG HELP					- reply with a helpful message
//
func (r *Handler) JSONDebug(subCmd rjs.DebugSubCommand, key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONDebug(subCmd, key, path)
}

// JSONForget is an alias for JSONDel
//
// ReJSON syntax:
// 	JSON.FORGET <key> [path]
//
func (r *Handler) JSONForget(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONForget(key, path)
}

// JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//
// ReJSON syntax:
// 	JSON.RESP <key> [path]
//
func (r *Handler) JSONResp(key, path string) (res interface{}, err error) {
	if r.clientName == rjs.ClientInactive {
		return nil, rjs.ErrNoClientSet
	}
	return r.implementation.JSONResp(key, path)
}
