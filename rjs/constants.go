package rjs

const (
	// ClientInactive signifies that the client is inactive in Handler
	ClientInactive = "inactive"

	// PopArrLast gives index of the last element for JSONArrPop
	PopArrLast = -1

	// DebugMemorySubcommand provide the corresponding MEMORY sub commands for JSONDebug
	DebugMemorySubcommand = "MEMORY"

	// DebugHelpSubcommand provide the corresponding HELP sub commands for JSONDebug
	DebugHelpSubcommand = "HELP"

	// DebugHelpOutput is the output of command JSON.Debug HELP <obj> [path]
	DebugHelpOutput = "MEMORY <key> [path] - reports memory usage\nHELP                - this message"

	// ReJSON Commands
	ReJSONCommand_SET       ReJSONCommandID = 0
	ReJSONCommand_GET       ReJSONCommandID = 1
	ReJSONCommand_DEL       ReJSONCommandID = 2
	ReJSONCommand_MGET      ReJSONCommandID = 3
	ReJSONCommand_TYPE      ReJSONCommandID = 4
	ReJSONCommand_NUMINCRBY ReJSONCommandID = 5
	ReJSONCommand_NUMMULTBY ReJSONCommandID = 6
	ReJSONCommand_STRAPPEND ReJSONCommandID = 7
	ReJSONCommand_STRLEN    ReJSONCommandID = 8
	ReJSONCommand_ARRAPPEND ReJSONCommandID = 9
	ReJSONCommand_ARRLEN    ReJSONCommandID = 10
	ReJSONCommand_ARRPOP    ReJSONCommandID = 11
	ReJSONCommand_ARRINDEX  ReJSONCommandID = 12
	ReJSONCommand_ARRTRIM   ReJSONCommandID = 13
	ReJSONCommand_ARRINSERT ReJSONCommandID = 14
	ReJSONCommand_OBJKEYS   ReJSONCommandID = 15
	ReJSONCommand_OBJLEN    ReJSONCommandID = 16
	ReJSONCommand_DEBUG     ReJSONCommandID = 17
	ReJSONCommand_FORGET    ReJSONCommandID = 18
	ReJSONCommand_RESP      ReJSONCommandID = 19

	// JSONGet Command Options
	GETOption_INDENT   GetOption = "INDENT"
	GETOption_NEWLINE  GetOption = "NEWLINE"
	GETOption_SPACE    GetOption = "SPACE"
	GETOption_NOESCAPE GetOption = "NOESCAPE"

	// JSONSET command Options
	SetOption_NX SetOption = "NX"
	SetOption_XX SetOption = "XX"
)

// commandName maps command id to the command name
var commandName = map[ReJSONCommandID]string{
	ReJSONCommand_SET:       "JSON.SET",
	ReJSONCommand_GET:       "JSON.GET",
	ReJSONCommand_DEL:       "JSON.DEL",
	ReJSONCommand_MGET:      "JSON.MGET",
	ReJSONCommand_TYPE:      "JSON.TYPE",
	ReJSONCommand_NUMINCRBY: "JSON.NUMINCRBY",
	ReJSONCommand_NUMMULTBY: "JSON.NUMMULTBY",
	ReJSONCommand_STRAPPEND: "JSON.STRAPPEND",
	ReJSONCommand_STRLEN:    "JSON.STRLEN",
	ReJSONCommand_ARRAPPEND: "JSON.ARRAPPEND",
	ReJSONCommand_ARRLEN:    "JSON.ARRLEN",
	ReJSONCommand_ARRPOP:    "JSON.ARRPOP",
	ReJSONCommand_ARRINDEX:  "JSON.ARRINDEX",
	ReJSONCommand_ARRTRIM:   "JSON.ARRTRIM",
	ReJSONCommand_ARRINSERT: "JSON.ARRINSERT",
	ReJSONCommand_OBJKEYS:   "JSON.OBJKEYS",
	ReJSONCommand_OBJLEN:    "JSON.OBJLEN",
	ReJSONCommand_DEBUG:     "JSON.DEBUG",
	ReJSONCommand_FORGET:    "JSON.FORGET",
	ReJSONCommand_RESP:      "JSON.RESP",
}

// commandMux maps command id to their Command Builder functions
var commandMux = map[ReJSONCommandID]CommandBuilderFunc{
	ReJSONCommand_SET:       commandJSONSet,
	ReJSONCommand_GET:       commandJSONGet,
	ReJSONCommand_DEL:       commandJSONgeneric,
	ReJSONCommand_MGET:      commandJSONMGet,
	ReJSONCommand_TYPE:      commandJSONgeneric,
	ReJSONCommand_NUMINCRBY: commandJSONNumIncrBy,
	ReJSONCommand_NUMMULTBY: commandJSONNumMultBy,
	ReJSONCommand_STRAPPEND: commandJSONStrAppend,
	ReJSONCommand_STRLEN:    commandJSONgeneric,
	ReJSONCommand_ARRAPPEND: commandJSONArrAppend,
	ReJSONCommand_ARRLEN:    commandJSONgeneric,
	ReJSONCommand_ARRPOP:    commandJSONArrPop,
	ReJSONCommand_ARRINDEX:  commandJSONArrIndex,
	ReJSONCommand_ARRTRIM:   commandJSONArrTrim,
	ReJSONCommand_ARRINSERT: commandJSONArrInsert,
	ReJSONCommand_OBJKEYS:   commandJSONgeneric,
	ReJSONCommand_OBJLEN:    commandJSONgeneric,
	ReJSONCommand_DEBUG:     commandJSONDebug,
	ReJSONCommand_FORGET:    commandJSONgeneric,
	ReJSONCommand_RESP:      commandJSONgeneric,
}
