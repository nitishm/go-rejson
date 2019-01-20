package rejson

/*
import (
	"fmt"
	"github.com/Shivam010/go-rejson/internal"
	"github.com/Shivam010/go-rejson/rjs"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	// PopArrLast gives index of the last element for JSONArrPop
	PopArrLast = -1
	// DebugMemorySubcommand provide the corresponding MEMORY sub commands for JSONDebug
	DebugMemorySubcommand = "MEMORY"
	// DebugHelpSubcommand provide the corresponding HELP sub commands for JSONDebug
	DebugHelpSubcommand = "HELP"
	// DebugHelpOutput is the output of command JSON.Debug HELP <obj> [path]
	DebugHelpOutput = "MEMORY <key> [path] - reports memory usage\nHELP                - this message"
)

// JSONGetOption provides methods various options for the JSON.GET Method
type JSONGetOption interface {
	optionTypeValue() (string, string)
}

// INDENT sets the indentation string for nested levels
type getOptionIndent struct {
	indentation string
}

func (opt *getOptionIndent) optionTypeValue() (string, string) {
	return "INDENT", opt.indentation
}

// NEWLINE sets the string that's printed at the end of each line
type getOptionNewLine struct {
	lineBreak string
}

func (opt *getOptionNewLine) optionTypeValue() (string, string) {
	return "NEWLINE", opt.lineBreak
}

// SPACE sets the string that's put between a key and a value
type getOptionSpace struct {
	space string
}

func (opt *getOptionSpace) optionTypeValue() (string, string) {
	return "SPACE", opt.space
}

// NOESCAPE option will disable the sending of \uXXXX escapes for non-ascii characters
type getOptionNoEscape struct{}

func (opt *getOptionNoEscape) optionTypeValue() (string, string) {
	return "NOESCAPE", ""
}

// NewJSONGetOptionIndent provides new INDENT options for JSON.GET Method
func NewJSONGetOptionIndent(val string) JSONGetOption { return &getOptionIndent{val} }

// NewJSONGetOptionNewLine provides new NEWLINE options for JSON.GET Method
func NewJSONGetOptionNewLine(val string) JSONGetOption { return &getOptionNewLine{val} }

// NewJSONGetOptionSpace provides new SPACE options for JSON.GET Method
func NewJSONGetOptionSpace(val string) JSONGetOption { return &getOptionSpace{val} }

// NewJSONGetOptionNoEscape provides new NOESCAPE options for JSON.GET Method
func NewJSONGetOptionNoEscape() JSONGetOption { return &getOptionNoEscape{} }

// JSONSet used to set a json object
// JSON.SET <key> <path> <json>
// 		 [NX | XX]
func JSONSet(conn redis.Conn, key string, path string, obj interface{}, NX bool, XX bool) (res interface{}, err error) {
	name, args, err := rjs.CommandBuilder("JSON.SET", key, path, obj, NX, XX)
	if err != nil {
		return nil, err
	}
	return conn.Do(name, args...)
}

// JSONGet used to get a json object
// JSON.GET <key>
//      [INDENT indentation-string]
// 		[NEWLINE line-break-string]
// 		[SPACE space-string]
// 		[NOESCAPE]
// 		[path ...]
func JSONGet(conn redis.Conn, key string, path string, opts ...JSONGetOption) (res interface{}, err error) {
	args := make([]interface{}, 0)
	args = append(args, key)

	for _, op := range opts {
		ty, va := op.optionTypeValue()

		args = append(args, ty)
		if ty != "NOESCAPE" {
			args = append(args, va)
		}
	}
	args = append(args, path)

	name, args, _ := rjs.CommandBuilder("JSON.GET", args...)
	return conn.Do(name, args...)
}

// JSONMGet used to get path values from multiple keys
// JSON.MGET <key> [key ...] <path>
func JSONMGet(conn redis.Conn, path string, keys ...string) (res interface{}, err error) {
	if len(keys) == 0 {
		err = fmt.Errorf("need atleast one key as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	args = append(args, path)
	name, args, _ := rjs.CommandBuilder("JSON.MGET", args...)
	return conn.Do(name, args...)
}

// JSONDel to delete a json object
// JSON.DEL <key> <path>
func JSONDel(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.DEL", key, path)
	return conn.Do(name, args...)
}

// JSONType to get the type of key or member at path.
// JSON.TYPE <key> [path]
func JSONType(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.TYPE", key, path)
	return conn.Do(name, args...)
}

// JSONNumIncrBy to increment a number by provided amount
// JSON.NUMINCRBY <key> <path> <number>
func JSONNumIncrBy(conn redis.Conn, key string, path string, number int) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.NUMINCRBY", key, path, number)
	return conn.Do(name, args...)
}

// JSONNumMultBy to increment a number by provided amount
// JSON.NUMMULTBY <key> <path> <number>
func JSONNumMultBy(conn redis.Conn, key string, path string, number int) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.NUMMULTBY", key, path, number)
	return conn.Do(name, args...)
}

// JSONStrAppend to append a jsonstring to an existing member
// JSON.STRAPPEND <key> [path] <json-string>
func JSONStrAppend(conn redis.Conn, key string, path string, jsonstring string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.STRAPPEND", key, path, jsonstring)
	return conn.Do(name, args...)
}

// JSONStrLen to return the length of a string member
// JSON.STRLEN <key> [path]
func JSONStrLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.STRLEN", key, path)
	return conn.Do(name, args...)
}

// JSONArrAppend to append json value into array at path
// JSON.ARRAPPEND <key> <path> <json> [json ...]
func JSONArrAppend(conn redis.Conn, key string, path string, values ...interface{}) (res interface{}, err error) {
	if len(values) == 0 {
		err = fmt.Errorf("need atleast one value string as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, key, path)
	args = append(args, values...)
	name, args, _ := rjs.CommandBuilder("JSON.ARRAPPEND", args...)
	return conn.Do(name, args...)
}

// JSONArrLen returns the length of the json array at path
// JSON.ARRLEN <key> [path]
func JSONArrLen(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.ARRLEN", key, path)
	return conn.Do(name, args...)
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
// JSON.ARRPOP <key> [path [index]]
func JSONArrPop(conn redis.Conn, key, path string, index int) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.ARRPOP", key, path, index)
	return conn.Do(name, args...)
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
// JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
func JSONArrIndex(conn redis.Conn, key, path string, jsonValue interface{}, optionalRange ...int) (
	res interface{}, err error) {

	args := []interface{}{key, path, jsonValue}

	ln := len(optionalRange)
	switch {
	case ln > 2:
		return nil, fmt.Errorf("need atmost two integeral value as an argument representing array range")
	case ln == 1: // only inclusive start is present
		args = append(args, optionalRange[0])
	case ln == 2: // both inclusive start and exclusive end are present
		args = append(args, optionalRange[0], optionalRange[1])
	}
	name, args, _ := rjs.CommandBuilder("JSON.ARRINDEX", args...)
	return conn.Do(name, args...)
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
// JSON.ARRTRIM <key> <path> <start> <stop>
func JSONArrTrim(conn redis.Conn, key, path string, start, end int) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.ARRTRIM", key, path, start, end)
	return conn.Do(name, args...)
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
// JSON.ARRINSERT <key> <path> <index> <json> [json ...]
func JSONArrInsert(conn redis.Conn, key, path string, index int, values ...interface{}) (res interface{}, err error) {
	if len(values) == 0 {
		err = fmt.Errorf("need atleast one value string as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, key, path, index)
	args = append(args, values...)
	name, args, _ := rjs.CommandBuilder("JSON.ARRINSERT", args...)
	return conn.Do(name, args...)
}

// JSONObjKeys returns the keys in the object that's referenced by path
// JSON.OBJKEYS <key> [path]
func JSONObjKeys(conn redis.Conn, key, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.OBJKEYS", key, path)
	res, err = conn.Do(name, args...)
	if err != nil {
		return
	}
	// JSON.OBJKEYS returns slice of string as slice of uint8
	slc := make([]string, 0, 10)
	for _, r := range res.([]interface{}) {
		slc = append(slc, rjs.BytesToString(r))
	}
	res = slc
	return
}

// JSONObjLen report the number of keys in the JSON Object at path in key
// JSON.OBJLEN <key> [path]
func JSONObjLen(conn redis.Conn, key, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.OBJLEN", key, path)
	return conn.Do(name, args...)
}

// JSONDebug reports information
// JSON.DEBUG <subcommand & arguments>
//		JSON.DEBUG MEMORY <key> [path]	- report the memory usage in bytes of a value. path defaults to root if not provided.
//		JSON.DEBUG HELP					- reply with a helpful message
func JSONDebug(conn redis.Conn, subcommand, key, path string) (res interface{}, err error) {
	if subcommand != DebugMemorySubcommand && subcommand != DebugHelpSubcommand {
		err = fmt.Errorf("unknown subcommand - try `JSON.DEBUG HELP`")
		return
	}
	name, args, _ := rjs.CommandBuilder("JSON.DEBUG", subcommand, key, path)
	res, err = conn.Do(name, args...)
	if err != nil {
		return
	}
	// JSONDebugMemorySubcommand returns an integer representing memory usage
	if subcommand == DebugMemorySubcommand {
		return res.(int64), err
	}
	// JSONDebugHelpSubcommand returns slice of string of Help as slice of uint8
	hlp := make([]string, 0, 10)
	for _, r := range res.([]interface{}) {
		hlp = append(hlp, rjs.BytesToString(r))
	}
	res = strings.Join(hlp, "\n")
	return
}

//JSONForget is an alias for JSONDel
//
func JSONForget(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.FORGET", key, path)
	return conn.Do(name, args...)
}

//JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//JSON.RESP <key> [path]
func JSONResp(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := rjs.CommandBuilder("JSON.RESP", key, path)
	return conn.Do(name, args...)
}
*/
