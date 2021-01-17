package rjs

import "fmt"

// Generic and Constant Errors to be returned to client to maintain stability
var (
	ErrInternal          = fmt.Errorf("error: internal client error")
	ErrNoClientSet       = fmt.Errorf("no redis client is set")
	ErrTooManyOptionals  = fmt.Errorf("error: too many optional arguments")
	ErrNeedAtLeastOneArg = fmt.Errorf("error: need atleast one argument in varying field")

	// GoRedis specific Nil error
	ErrGoRedisNil = fmt.Errorf("redis: nil")
)

const (
	// ClientInactive signifies that the client is inactive in Handler
	ClientInactive = "inactive"

	// ClientRedigo signifies that the current client is redigo
	ClientRedigo = "redigo"

	// ClientGoRedis signifies that the current client is go-redis
	ClientGoRedis = "goredis"

	// PopArrLast gives index of the last element for JSONArrPop
	PopArrLast = -1

	// DebugMemorySubcommand provide the corresponding MEMORY sub commands for JSONDebug
	DebugMemorySubcommand DebugSubCommand = "MEMORY"

	// DebugHelpSubcommand provide the corresponding HELP sub commands for JSONDebug
	DebugHelpSubcommand DebugSubCommand = "HELP"

	// DebugHelpOutput is the output of command JSON.Debug HELP <obj> [path]
	DebugHelpOutput = "MEMORY <key> [path] - reports memory usage\nHELP                - this message"

	// ReJSON Commands
	ReJSONCommandSET       ReJSONCommandID = 0
	ReJSONCommandGET       ReJSONCommandID = 1
	ReJSONCommandDEL       ReJSONCommandID = 2
	ReJSONCommandMGET      ReJSONCommandID = 3
	ReJSONCommandTYPE      ReJSONCommandID = 4
	ReJSONCommandNUMINCRBY ReJSONCommandID = 5
	ReJSONCommandNUMMULTBY ReJSONCommandID = 6
	ReJSONCommandSTRAPPEND ReJSONCommandID = 7
	ReJSONCommandSTRLEN    ReJSONCommandID = 8
	ReJSONCommandARRAPPEND ReJSONCommandID = 9
	ReJSONCommandARRLEN    ReJSONCommandID = 10
	ReJSONCommandARRPOP    ReJSONCommandID = 11
	ReJSONCommandARRINDEX  ReJSONCommandID = 12
	ReJSONCommandARRTRIM   ReJSONCommandID = 13
	ReJSONCommandARRINSERT ReJSONCommandID = 14
	ReJSONCommandOBJKEYS   ReJSONCommandID = 15
	ReJSONCommandOBJLEN    ReJSONCommandID = 16
	ReJSONCommandDEBUG     ReJSONCommandID = 17
	ReJSONCommandFORGET    ReJSONCommandID = 18
	ReJSONCommandRESP      ReJSONCommandID = 19

	// JSONSET command Options
	SetOptionNX SetOption = "NX"
	SetOptionXX SetOption = "XX"
)

// JSONGet Command Options
var (
	GETOptionSPACE    = GetOption{"SPACE", " "}
	GETOptionINDENT   = GetOption{"INDENT", "\t"}
	GETOptionNEWLINE  = GetOption{"NEWLINE", "\n"}
	GETOptionNOESCAPE = GetOption{"NOESCAPE", ""}
)

// commandName maps command id to the command name
var commandName = map[ReJSONCommandID]string{
	ReJSONCommandSET:       "JSON.SET",
	ReJSONCommandGET:       "JSON.GET",
	ReJSONCommandDEL:       "JSON.DEL",
	ReJSONCommandMGET:      "JSON.MGET",
	ReJSONCommandTYPE:      "JSON.TYPE",
	ReJSONCommandNUMINCRBY: "JSON.NUMINCRBY",
	ReJSONCommandNUMMULTBY: "JSON.NUMMULTBY",
	ReJSONCommandSTRAPPEND: "JSON.STRAPPEND",
	ReJSONCommandSTRLEN:    "JSON.STRLEN",
	ReJSONCommandARRAPPEND: "JSON.ARRAPPEND",
	ReJSONCommandARRLEN:    "JSON.ARRLEN",
	ReJSONCommandARRPOP:    "JSON.ARRPOP",
	ReJSONCommandARRINDEX:  "JSON.ARRINDEX",
	ReJSONCommandARRTRIM:   "JSON.ARRTRIM",
	ReJSONCommandARRINSERT: "JSON.ARRINSERT",
	ReJSONCommandOBJKEYS:   "JSON.OBJKEYS",
	ReJSONCommandOBJLEN:    "JSON.OBJLEN",
	ReJSONCommandDEBUG:     "JSON.DEBUG",
	ReJSONCommandFORGET:    "JSON.FORGET",
	ReJSONCommandRESP:      "JSON.RESP",
}

// commandMux maps command id to their Command Builder functions
var commandMux = map[ReJSONCommandID]CommandBuilderFunc{
	ReJSONCommandSET:       commandJSONSet,
	ReJSONCommandGET:       commandJSONGet,
	ReJSONCommandDEL:       commandJSONGeneric,
	ReJSONCommandMGET:      commandJSONMGet,
	ReJSONCommandTYPE:      commandJSONGeneric,
	ReJSONCommandNUMINCRBY: commandJSONNumIncrBy,
	ReJSONCommandNUMMULTBY: commandJSONNumMultBy,
	ReJSONCommandSTRAPPEND: commandJSONStrAppend,
	ReJSONCommandSTRLEN:    commandJSONGeneric,
	ReJSONCommandARRAPPEND: commandJSONArrAppend,
	ReJSONCommandARRLEN:    commandJSONGeneric,
	ReJSONCommandARRPOP:    commandJSONArrPop,
	ReJSONCommandARRINDEX:  commandJSONArrIndex,
	ReJSONCommandARRTRIM:   commandJSONArrTrim,
	ReJSONCommandARRINSERT: commandJSONArrInsert,
	ReJSONCommandOBJKEYS:   commandJSONGeneric,
	ReJSONCommandOBJLEN:    commandJSONGeneric,
	ReJSONCommandDEBUG:     commandJSONDebug,
	ReJSONCommandFORGET:    commandJSONGeneric,
	ReJSONCommandRESP:      commandJSONGeneric,
}
