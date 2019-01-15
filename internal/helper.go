package rjs

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

