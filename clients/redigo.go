package clients

import (
	"fmt"
	"github.com/nitishm/go-rejson/v4/rjs"
	"strings"
)

// RedigoClientConn - an abstracted interface for redigo.Conn and redigo.ConnWithTimeout
type RedigoClientConn interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

// Redigo implements ReJSON interface for GoModule/Redigo Redis client
// Link: https://github.com/gomodule/redigo
type Redigo struct {
	Conn RedigoClientConn
}

// JSONSet used to set a json object
//
// ReJSON syntax:
// 	JSON.SET <key> <path> <json>
// 			 [NX | XX]
//
func (r *Redigo) JSONSet(key string, path string, obj interface{}, opts ...rjs.SetOption) (res interface{}, err error) {

	if len(opts) > 1 {
		return nil, rjs.ErrTooManyOptionals
	}
	args := make([]interface{}, 0, 5)
	args = append(args, key, path, obj)

	if len(opts) == 1 {
		args = append(args, opts[0].Value()...)
	}
	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandSET, args...)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
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
func (r *Redigo) JSONGet(key, path string, opts ...rjs.GetOption) (res interface{}, err error) {

	if len(opts) > 4 {
		return nil, rjs.ErrTooManyOptionals
	}
	args := make([]interface{}, 0)
	args = append(args, key, path)

	for _, op := range opts {
		args = append(args, op.Value()...)
	}

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandGET, args...)
	if err != nil {
		return nil, err
	}

	return r.Conn.Do(name, args...)
}

// JSONMGet used to get path values from multiple keys
//
// ReJSON syntax:
// 	JSON.MGET <key> [key ...] <path>
//
func (r *Redigo) JSONMGet(path string, keys ...string) (res interface{}, err error) {

	if len(keys) == 0 {
		return nil, rjs.ErrNeedAtLeastOneArg
	}
	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}

	args = append(args, path)
	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandMGET, args...)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONDel used to delete a json object
//
// ReJSON syntax:
// 	JSON.DEL <key> <path>
//
func (r *Redigo) JSONDel(key string, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandDEL, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONType used to get the type of key or member at path.
//
// ReJSON syntax:
// 	JSON.TYPE <key> [path]
//
func (r *Redigo) JSONType(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandTYPE, key, path)
	if err != nil {
		return nil, err
	}

	res, err = r.Conn.Do(name, args...)

	if err != nil {
		return nil, err
	}
	switch v := res.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case nil:
		return
	default:
		err := fmt.Errorf("type returned not expected %T", v)
		return nil, err
	}
}

// JSONNumIncrBy used to increment a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMINCRBY <key> <path> <number>
//
func (r *Redigo) JSONNumIncrBy(key, path string, number int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandNUMINCRBY, key, path, number)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONNumMultBy to multiply a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMMULTBY <key> <path> <number>
//
func (r *Redigo) JSONNumMultBy(key, path string, number int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandNUMMULTBY, key, path, number)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONStrAppend used to append a jsonstring to an existing member
//
// ReJSON syntax:
// 	JSON.STRAPPEND <key> [path] <json-string>
//
func (r *Redigo) JSONStrAppend(key, path, jsonstring string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandSTRAPPEND, key, path, jsonstring)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONStrLen used to return the length of a string member
//
// ReJSON syntax:
// 	JSON.STRLEN <key> [path]
//
func (r *Redigo) JSONStrLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandSTRLEN, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrAppend used to append json value into array at path
//
// ReJSON syntax:
// 	JSON.ARRAPPEND <key> <path> <json> [json ...]
//
func (r *Redigo) JSONArrAppend(key, path string, values ...interface{}) (res interface{}, err error) {

	if len(values) == 0 {
		return nil, rjs.ErrNeedAtLeastOneArg
	}
	args := make([]interface{}, 0)
	args = append(args, key, path)
	args = append(args, values...)

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRAPPEND, args...)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrLen returns the length of the json array at path
//
// ReJSON syntax:
// 	JSON.ARRLEN <key> [path]
//
func (r *Redigo) JSONArrLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRLEN, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
//
// ReJSON syntax:
// 	JSON.ARRPOP <key> [path [index]]
//
func (r *Redigo) JSONArrPop(key, path string, index int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRPOP, key, path, index)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
//
// ReJSON syntax:
// 	JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
//
func (r *Redigo) JSONArrIndex(key, path string, jsonValue interface{}, optionalRange ...int) (res interface{}, err error) { // nolint: lll

	args := []interface{}{key, path, jsonValue}

	ln := len(optionalRange)
	switch {
	case ln > 2:
		return nil, rjs.ErrTooManyOptionals
	case ln == 1: // only inclusive start is present
		args = append(args, optionalRange[0])
	case ln == 2: // both inclusive start and exclusive end are present
		args = append(args, optionalRange[0], optionalRange[1])
	}
	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRINDEX, args...)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
//
// ReJSON syntax:
// 	JSON.ARRTRIM <key> <path> <start> <stop>
//
func (r *Redigo) JSONArrTrim(key, path string, start, end int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRTRIM, key, path, start, end)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
//
// ReJSON syntax:
// 	JSON.ARRINSERT <key> <path> <index> <json> [json ...]
//
func (r *Redigo) JSONArrInsert(key, path string, index int, values ...interface{}) (res interface{}, err error) {

	if len(values) == 0 {
		return nil, rjs.ErrNeedAtLeastOneArg
	}
	args := make([]interface{}, 0)
	args = append(args, key, path, index)
	args = append(args, values...)

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRINSERT, args...)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONObjKeys returns the keys in the object that's referenced by path
//
// ReJSON syntax:
// 	JSON.OBJKEYS <key> [path]
//
func (r *Redigo) JSONObjKeys(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandOBJKEYS, key, path)
	if err != nil {
		return nil, err
	}
	res, err = r.Conn.Do(name, args...)
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
//
// ReJSON syntax:
// 	JSON.OBJLEN <key> [path]
//
func (r *Redigo) JSONObjLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandOBJLEN, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONDebug reports information
//
// ReJSON syntax:
// 	JSON.DEBUG <subcommand & arguments>
//		JSON.DEBUG MEMORY <key> [path]	- report the memory usage in bytes of a value. path defaults to root if not provided.
//		JSON.DEBUG HELP					- reply with a helpful message
//
func (r *Redigo) JSONDebug(subcommand rjs.DebugSubCommand, key, path string) (res interface{}, err error) {

	if subcommand != rjs.DebugMemorySubcommand && subcommand != rjs.DebugHelpSubcommand {
		err = fmt.Errorf("unknown subcommand - try `JSON.DEBUG HELP`")
		return
	}
	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandDEBUG, subcommand, key, path)
	if err != nil {
		return nil, err
	}
	res, err = r.Conn.Do(name, args...)
	if err != nil {
		return
	}
	// JSONDebugMemorySubcommand returns an integer representing memory usage
	if subcommand == rjs.DebugMemorySubcommand {
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

// JSONForget is an alias for JSONDel
//
// ReJSON syntax:
// 	JSON.FORGET <key> [path]
//
func (r *Redigo) JSONForget(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandFORGET, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}

// JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//
// ReJSON syntax:
// 	JSON.RESP <key> [path]
//
func (r *Redigo) JSONResp(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandRESP, key, path)
	if err != nil {
		return nil, err
	}
	return r.Conn.Do(name, args...)
}
