package rejson

import (
	"encoding/json"
	"fmt"
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
	// DebugHelpOutput is the ouput of command JSON.Debug HELP <obj> [path]
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

// commandMux maps command name to a command function
var commandMux = map[string]func(argsIn ...interface{}) (argsOut []interface{}, err error){
	"JSON.SET":       commandJSONSet,
	"JSON.GET":       commandJSONGet,
	"JSON.DEL":       commandJSON,
	"JSON.MGET":      commandJSONMGet,
	"JSON.TYPE":      commandJSON,
	"JSON.NUMINCRBY": commandJSONNumIncrBy,
	"JSON.NUMMULTBY": commandJSONNumMultBy,
	"JSON.STRAPPEND": commandJSONStrAppend,
	"JSON.STRLEN":    commandJSON,
	"JSON.ARRAPPEND": commandJSONArrAppend,
	"JSON.ARRLEN":    commandJSON,
	"JSON.ARRPOP":    commandJSONArrPop,
	"JSON.ARRINDEX":  commandJSONArrIndex,
	"JSON.ARRTRIM":   commandJSONArrTrim,
	"JSON.ARRINSERT": commandJSONArrInsert,
	"JSON.OBJKEYS":   commandJSON,
	"JSON.OBJLEN":    commandJSON,
	"JSON.DEBUG":     commandJSONDebug,
	"JSON.FORGET":    commandJSON,
	"JSON.RESP":      commandJSON,
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
		err = fmt.Errorf("both NX and XX cannot be true")
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
	argsOut = append(argsOut, argsIn[2:]...)
	return
}

func commandJSON(argsIn ...interface{}) (argsOut []interface{}, err error) {
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

func commandJSONDebug(argsIn ...interface{}) (argsOut []interface{}, err error) {
	subcommand := argsIn[0]
	key := argsIn[1]
	path := argsIn[2]
	argsOut = append(argsOut, subcommand, key, path)
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

	name, args, _ := CommandBuilder("JSON.GET", args...)
	return conn.Do(name, args...)
}

// JSONMGet used to get path values from multiple keys
// JSON.MGET <key> [key ...] <path>
func JSONMGet(conn redis.Conn, path string, keys ...string) (res interface{}, err error) {
	if len(keys) == 0 {
		err = fmt.Errorf("need atlesat one key as an argument")
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
		err = fmt.Errorf("need atlesat one value string as an argument")
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
		return nil, fmt.Errorf("need atmost two integeral value as an argument representing array range")
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
		err = fmt.Errorf("need atlesat one value string as an argument")
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, key, path, index)
	args = append(args, values...)
	name, args, _ := CommandBuilder("JSON.ARRINSERT", args...)
	return conn.Do(name, args...)
}

// JSONObjKeys returns the keys in the object that's referenced by path
// JSON.OBJKEYS <key> [path]
func JSONObjKeys(conn redis.Conn, key, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.OBJKEYS", key, path)
	res, err = conn.Do(name, args...)
	if err != nil {
		return
	}
	// JSON.OBJKEYS returns slice of string as slice of uint8
	slc := make([]string, 0, 10)
	for _, r := range res.([]interface{}) {
		slc = append(slc, tostring(r))
	}
	res = slc
	return
}

// JSONObjLen report the number of keys in the JSON Object at path in key
// JSON.OBJLEN <key> [path]
func JSONObjLen(conn redis.Conn, key, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.OBJLEN", key, path)
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
	name, args, _ := CommandBuilder("JSON.DEBUG", subcommand, key, path)
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
		hlp = append(hlp, tostring(r))
	}
	res = strings.Join(hlp, "\n")
	return
}

//JSONForget is an alias for JSONDel
func JSONForget(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.FORGET", key, path)
	return conn.Do(name, args...)
}

//JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//JSON.RESP <key> [path]
func JSONResp(conn redis.Conn, key string, path string) (res interface{}, err error) {
	name, args, _ := CommandBuilder("JSON.RESP", key, path)
	return conn.Do(name, args...)
}

// tostring converts each byte in slice into character, else panic out
func tostring(lst interface{}) (str string) {
	_lst, ok := lst.([]byte)
	if !ok {
		panic("error: something went wrong")
	}
	for _, s := range _lst {
		str += string(s)
	}
	return
}
