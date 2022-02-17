package clients

import (
	"context"
	"fmt"
	"strings"

	goredis "github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4/rjs"
)

// GoRedisClientConn - an abstracted interface for goredis.Client, goredis.ClusterClient, goredis.Ring,
// or goredis.UniversalClient
type GoRedisClientConn interface {
	Do(ctx context.Context, args ...interface{}) *goredis.Cmd
}

// GoRedis implements ReJSON interface for Go-Redis/Redis Redis client
// Link: https://github.com/go-redis/redis
type GoRedis struct {
	Conn GoRedisClientConn
	// ctx defines context for the provided connection
	ctx context.Context
}

// NewGoRedisClient returns a new GoRedis ReJSON client with the provided context
// and connection, if ctx is nil default context.Background will be used
func NewGoRedisClient(ctx context.Context, conn GoRedisClientConn) *GoRedis {
	if ctx == nil {
		ctx = context.Background()
	}
	return &GoRedis{
		ctx:  ctx,
		Conn: conn,
	}
}

// JSONSet used to set a json object
//
// ReJSON syntax:
// 	JSON.SET <key> <path> <json>
// 			 [NX | XX]
//
func (r *GoRedis) JSONSet(key string, path string, obj interface{}, opts ...rjs.SetOption) (res interface{}, err error) { // nolint: lll

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
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()

	if err != nil && err.Error() == rjs.ErrGoRedisNil.Error() {
		err = nil
	}
	return
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
func (r *GoRedis) JSONGet(key, path string, opts ...rjs.GetOption) (res interface{}, err error) {

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

	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	return rjs.StringToBytes(res), err
}

// JSONMGet used to get path values from multiple keys
//
// ReJSON syntax:
// 	JSON.MGET <key> [key ...] <path>
//
func (r *GoRedis) JSONMGet(path string, keys ...string) (res interface{}, err error) {

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
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	nres := make([]interface{}, 0, 10)
	for _, r := range res.([]interface{}) {
		if r != nil {
			nres = append(nres, rjs.StringToBytes(r))
		} else {
			nres = append(nres, nil)
		}
	}

	return nres, nil
}

// JSONDel to delete a json object
//
// ReJSON syntax:
// 	JSON.DEL <key> <path>
//
func (r *GoRedis) JSONDel(key string, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandDEL, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONType to get the type of key or member at path.
//
// ReJSON syntax:
// 	JSON.TYPE <key> [path]
//
func (r *GoRedis) JSONType(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandTYPE, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()

	if err != nil && err.Error() == rjs.ErrGoRedisNil.Error() {
		err = nil
	}
	return
}

// JSONNumIncrBy to increment a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMINCRBY <key> <path> <number>
//
func (r *GoRedis) JSONNumIncrBy(key, path string, number int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandNUMINCRBY, key, path, number)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	return rjs.StringToBytes(res), err
}

// JSONNumMultBy to increment a number by provided amount
//
// ReJSON syntax:
// 	JSON.NUMMULTBY <key> <path> <number>
//
func (r *GoRedis) JSONNumMultBy(key, path string, number int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandNUMMULTBY, key, path, number)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	return rjs.StringToBytes(res), err
}

// JSONStrAppend to append a jsonstring to an existing member
//
// ReJSON syntax:
// 	JSON.STRAPPEND <key> [path] <json-string>
//
func (r *GoRedis) JSONStrAppend(key, path, jsonstring string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandSTRAPPEND, key, path, jsonstring)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONStrLen to return the length of a string member
//
// ReJSON syntax:
// 	JSON.STRLEN <key> [path]
//
func (r *GoRedis) JSONStrLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandSTRLEN, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONArrAppend to append json value into array at path
//
// ReJSON syntax:
// 	JSON.ARRAPPEND <key> <path> <json> [json ...]
//
func (r *GoRedis) JSONArrAppend(key, path string, values ...interface{}) (res interface{}, err error) {

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
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONArrLen returns the length of the json array at path
//
// ReJSON syntax:
// 	JSON.ARRLEN <key> [path]
//
func (r *GoRedis) JSONArrLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRLEN, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONArrPop removes and returns element from the index in the array
// to pop last element use rejson.PopArrLast
//
// ReJSON syntax:
// 	JSON.ARRPOP <key> [path [index]]
//
func (r *GoRedis) JSONArrPop(key, path string, index int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRPOP, key, path, index)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)

	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	return rjs.StringToBytes(res), err
}

// JSONArrIndex returns the index of the json element provided and return -1 if element is not present
//
// ReJSON syntax:
// 	JSON.ARRINDEX <key> <path> <json-scalar> [start [stop]]
//
func (r *GoRedis) JSONArrIndex(key, path string, jsonValue interface{}, optionalRange ...int) (res interface{}, err error) { // nolint: lll

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
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
//
// ReJSON syntax:
// 	JSON.ARRTRIM <key> <path> <start> <stop>
//
func (r *GoRedis) JSONArrTrim(key, path string, start, end int) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandARRTRIM, key, path, start, end)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONArrInsert inserts the json value(s) into the array at path before the index (shifts to the right).
//
// ReJSON syntax:
// 	JSON.ARRINSERT <key> <path> <index> <json> [json ...]
//
func (r *GoRedis) JSONArrInsert(key, path string, index int, values ...interface{}) (res interface{}, err error) {

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
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONObjKeys returns the keys in the object that's referenced by path
//
// ReJSON syntax:
// 	JSON.OBJKEYS <key> [path]
//
func (r *GoRedis) JSONObjKeys(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandOBJKEYS, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
	if err != nil {
		return
	}
	// JSON.OBJKEYS returns slice of string as slice of uint8
	slc := make([]string, 0, 10)
	for _, r := range res.([]interface{}) {
		slc = append(slc, r.(string))
	}
	res = slc
	return
}

// JSONObjLen report the number of keys in the JSON Object at path in key
//
// ReJSON syntax:
// 	JSON.OBJLEN <key> [path]
//
func (r *GoRedis) JSONObjLen(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandOBJLEN, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONDebug reports information
//
// ReJSON syntax:
// 	JSON.DEBUG <subcommand & arguments>
//		JSON.DEBUG MEMORY <key> [path]	- report the memory usage in bytes of a value. path defaults to root if not provided.
//		JSON.DEBUG HELP					- reply with a helpful message
//
func (r *GoRedis) JSONDebug(subcommand rjs.DebugSubCommand, key, path string) (res interface{}, err error) {

	if subcommand != rjs.DebugMemorySubcommand && subcommand != rjs.DebugHelpSubcommand {
		err = fmt.Errorf("unknown subcommand - try `JSON.DEBUG HELP`")
		return
	}
	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandDEBUG, string(subcommand), key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	res, err = r.Conn.Do(r.ctx, args...).Result()
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
		hlp = append(hlp, r.(string))
	}
	res = strings.Join(hlp, "\n")
	return
}

// JSONForget is an alias for JSONDel
//
// ReJSON syntax:
// 	JSON.FORGET <key> [path]
//
func (r *GoRedis) JSONForget(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandFORGET, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}

// JSONResp returns the JSON in key in Redis Serialization Protocol (RESP).
//
// ReJSON syntax:
// 	JSON.RESP <key> [path]
//
func (r *GoRedis) JSONResp(key, path string) (res interface{}, err error) {

	name, args, err := rjs.CommandBuilder(rjs.ReJSONCommandRESP, key, path)
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{name}, args...)
	return r.Conn.Do(r.ctx, args...).Result()
}
