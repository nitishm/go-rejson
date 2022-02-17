package rjs

// ReJSONOption provides methods for the options used by various ReJSON commands
// It also abstracts options from the required parameters of the commands
//
// Like:
// 	JSON.GET, JSON.SET, JSON.ARRINDEX, JSON.ARRPOP
type ReJSONOption interface {
	// Value returns the value of the option being used
	Value() []interface{}

	// MethodID returns the ID of the ReJSON Function defined in ReJSONCommands
	// whose options are begins implemented
	MethodID() ReJSONCommandID
}

// GetOption implements ReJSONOption for JSON.GET Method
// Get Options:
// 	* INDENT 	(with default set to a tab, '\t')
// 	* NEWLINE	(with default set to a new line, '\n')
//  * SPACE		(with default set to a space, ' ')
//  * NOESCAPE  (a boolean type option)
type GetOption struct {
	name string
	Arg  string
}

// MethodID returns the name of the method i.e. JSON.GET
func (g GetOption) MethodID() ReJSONCommandID {
	return ReJSONCommandGET
}

// Value returns the value of the option being used
func (g GetOption) Value() []interface{} {
	if g.name == GETOptionNOESCAPE.name {
		return []interface{}{g.name}
	}
	return []interface{}{g.name, g.Arg}
}

// SetValue will set the values in the options
func (g *GetOption) SetValue(arg string) {
	g.Arg = arg
}

// SetOption implements ReJSONOption for JSON.SET Method
// Set Options:
//	* NX or XX
type SetOption string

// MethodID returns the name of the method i.e. JSON.SET
func (g SetOption) MethodID() ReJSONCommandID {
	return ReJSONCommandSET
}

// Value returns the value of the option being used
func (g SetOption) Value() []interface{} {
	return []interface{}{string(g)}
}
