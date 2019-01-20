package rjs

import "fmt"

// BytesToString converts each byte in a byte slice into character, else panic out
func BytesToString(lst interface{}) (str string) {
	_lst, ok := lst.([]byte)
	if !ok {
		panic("error: something went wrong")
	}
	for _, s := range _lst {
		str += string(s)
	}
	return
}

// Value returns integral value of the ReJSON Command ID
func (r ReJSONCommandId) Value() int32 {
	return int32(r)
}

// TypeSafety checks the validity of the command id
func (r ReJSONCommandId) TypeSafety() error {
	if r.Value() < 0 || r.Value() > 19 {
		return fmt.Errorf("error: invalid command id")
	}
	return nil
}

// Details returns the details of the CommandId like its command function and name
func (r ReJSONCommandId) Details() (CommandBuilderFunc, string, error) {
	name, ok := commandName[r]
	if !ok {
		return nil, "", fmt.Errorf("command not supported by ReJSON")
	}
	cmd := commandMux[r]
	return cmd, name, nil
}
