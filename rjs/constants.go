package rjs

const (
	// PopArrLast gives index of the last element for JSONArrPop
	PopArrLast = -1

	// DebugMemorySubcommand provide the corresponding MEMORY sub commands for JSONDebug
	DebugMemorySubcommand = "MEMORY"

	// DebugHelpSubcommand provide the corresponding HELP sub commands for JSONDebug
	DebugHelpSubcommand = "HELP"

	// DebugHelpOutput is the output of command JSON.Debug HELP <obj> [path]
	DebugHelpOutput = "MEMORY <key> [path] - reports memory usage\nHELP                - this message"

	// ReJSON Commands
	ReJSONCommand_SET       ReJSONCommandId = 0
	ReJSONCommand_GET       ReJSONCommandId = 1
	ReJSONCommand_DEL       ReJSONCommandId = 2
	ReJSONCommand_MGET      ReJSONCommandId = 3
	ReJSONCommand_TYPE      ReJSONCommandId = 4
	ReJSONCommand_NUMINCRBY ReJSONCommandId = 5
	ReJSONCommand_NUMMULTBY ReJSONCommandId = 6
	ReJSONCommand_STRAPPEND ReJSONCommandId = 7
	ReJSONCommand_STRLEN    ReJSONCommandId = 8
	ReJSONCommand_ARRAPPEND ReJSONCommandId = 9
	ReJSONCommand_ARRLEN    ReJSONCommandId = 10
	ReJSONCommand_ARRPOP    ReJSONCommandId = 11
	ReJSONCommand_ARRINDEX  ReJSONCommandId = 12
	ReJSONCommand_ARRTRIM   ReJSONCommandId = 13
	ReJSONCommand_ARRINSERT ReJSONCommandId = 14
	ReJSONCommand_OBJKEYS   ReJSONCommandId = 15
	ReJSONCommand_OBJLEN    ReJSONCommandId = 16
	ReJSONCommand_DEBUG     ReJSONCommandId = 17
	ReJSONCommand_FORGET    ReJSONCommandId = 18
	ReJSONCommand_RESP      ReJSONCommandId = 19

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
var commandName = map[ReJSONCommandId]string{
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
var commandMux = map[ReJSONCommandId]CommandBuilderFunc{
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
