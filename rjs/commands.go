package rjs

import (
	"encoding/json"
	"fmt"
)

// ReJSONCommandID marks a particular unique id to all the ReJSON commands
// to ensure proper type safety and help reducing typos in using them.
type ReJSONCommandID int32

// DebugSubCommand provides the abstract sub-commands for the JSON.DEBUG command
type DebugSubCommand string

// CommandBuilderFunc uses for the simplicity of the corresponding ReJSON module command builders
type CommandBuilderFunc func(argsIn ...interface{}) (argsOut []interface{}, err error)

func commandJSONSet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	obj := argsIn[2]

	argsOut = append(argsOut, key, path)

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	argsOut = append(argsOut, b)

	if len(argsIn) == 4 {
		argsOut = append(argsOut, argsIn[3])
	}
	return
}

func commandJSONGet(argsIn ...interface{}) (argsOut []interface{}, err error) {
	key := argsIn[0]
	path := argsIn[1]
	argsOut = append(argsOut, key)
	argsOut = append(argsOut, argsIn[2:]...)
	argsOut = append(argsOut, path)
	return
}

func commandJSONGeneric(argsIn ...interface{}) (argsOut []interface{}, err error) {
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
	// if index is not used as option for PopArrLast ( == -1 ), append index
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

// CommandBuilder is used to build a command that can be used directly to
// build REJSON commands
//
// This is especially useful if you do not need to client's `Do()` method
// and instead need to use the JSON.* commands in the MUTLI/EXEC scenario
// along with some other operations like
// 	GET/SET/HGET/HSET/...
//
func CommandBuilder(commandNameIn ReJSONCommandID, argsIn ...interface{}) (commandNameOut string, argsOut []interface{}, err error) { // nolint: lll

	cmd, commandNameOut, err := commandNameIn.Details()
	if err != nil {
		return
	}

	argsOut, err = cmd(argsIn...)
	if err != nil {
		return commandNameOut, nil, fmt.Errorf("failed to execute command %v: %v", commandNameIn, err)
	}

	return
}
