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

// StringToBytes converts each character of the string slice into byte, else panic out
func StringToBytes(lst interface{}) (by []byte) {
	_lst, ok := lst.(string)
	if !ok {
		panic("error: something went wrong")
	}
	by = []byte(_lst)
	return
}

// Value returns integral value of the ReJSON Command ID
func (r ReJSONCommandID) Value() int32 {
	return int32(r)
}

// TypeSafety checks the validity of the command id
func (r ReJSONCommandID) TypeSafety() error {
	if r.Value() < 0 || r.Value() > 19 {
		return fmt.Errorf("error: invalid command id")
	}
	return nil
}

// Details returns the details of the CommandId like its command function and name
func (r ReJSONCommandID) Details() (CommandBuilderFunc, string, error) {
	name, ok := commandName[r]
	if !ok {
		return nil, "", fmt.Errorf("command not supported by ReJSON")
	}
	cmd := commandMux[r]
	return cmd, name, nil
}
