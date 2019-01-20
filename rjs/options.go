package rjs

import "fmt"

// ReJSONOption provides methods for the options used by various ReJSON commands
//
// Like:
// JSON.Get, JSON.Set, etc
type ReJSONOption interface {
	// Value returns the value of the option being used
	Value() string

	// MethodName returns the name of the method whose options are begins implemented
	MethodName() string

	// UseOption is used to apply the option in the command
	//
	//
	// Examples:
	//
	//
	//
	//
	UseOption(...interface{}) ([]interface{}, error)

}

// GetOption implements ReJSONOption for JSON.GET Method
// Get Options:
// 	* INDENT
// 	* NEWLINE
//  * SPACE
//  * NOESCAPE
type GetOption string

// MethodName returns the name of the method i.e. JSON.GET
func (g GetOption) MethodName() string {
	return "JSON.GET"
}

// Value returns the value of the option being used
func (g GetOption) Value() string {
	return string(g)
}

// UseOption is used to apply the option in the command
func (g GetOption) UseOption(args ...interface{}) ([]interface{}, error) {

	if len(args) > 1 || (g == GETOption_NOESCAPE && len(args) != 0) {
		return nil, fmt.Errorf("error: too many arguments")
	}

	return []interface{}{
		g.Value(),
		args[0],
	}, nil
}

type SetOption string

// MethodName returns the name of the method i.e. JSON.SET
func (g SetOption) MethodName() string {
	return "JSON.SET"
}

// Value returns the value of the option being used
func (g SetOption) Value() string {
	return string(g)
}

// UseOption is used to apply the option in the command
func (g SetOption) UseOption(args ...interface{}) ([]interface{}, error) {

	if len(args) > 0 {
		return nil, fmt.Errorf("error: too many arguments")
	}

	return []interface{}{
		g.Value(),
	}, nil
}
