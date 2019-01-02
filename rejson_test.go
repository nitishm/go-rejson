package rejson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gomodule/redigo/redis"
)

type TestObject struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func TestJSONSet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	testObj := TestObject{
		"item#1",
		1,
	}
	type args struct {
		key  string
		path string
		obj  interface{}
		NX   bool
		XX   bool
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				key:  "kstr",
				path: ".",
				obj:  "simplestring",
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleInt",
			args: args{
				key:  "kint",
				path: ".",
				obj:  1234,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "kstruct",
				path: ".",
				obj:  testObj,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStructFieldOK",
			args: args{
				key:  "kstruct",
				path: "name",
				obj:  "foobar",
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStringWithNXOK",
			args: args{
				key:  "kstrnx",
				path: ".",
				obj:  123,
				NX:   true,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStringWithNXNotOK",
			args: args{
				key:  "kstrnx",
				path: ".",
				obj:  "simplestringnx",
				NX:   true,
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXOK",
			args: args{
				key:  "kstrnx",
				path: ".",
				obj:  "simplestringfoo",
				XX:   true,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXNotOK",
			args: args{
				key:  "kstrxx",
				path: ".",
				obj:  "simplestringfoobar",
				XX:   true,
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXNX",
			args: args{
				key:  "kstrxxnx",
				path: ".",
				obj:  "simplestringfoobar",
				XX:   true,
				NX:   true,
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONSet(conn, tt.args.key, tt.args.path, tt.args.obj, tt.args.NX, tt.args.XX)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONGet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kint", ".", 123, false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "kstruct", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				conn: conn,
				key:  "kstr",
				path: ".",
			},
			wantRes: []byte("\"simplestring\""),
			wantErr: false,
		},
		{
			name: "SimpleInt",
			args: args{
				conn: conn,
				key:  "kint",
				path: ".",
			},
			wantRes: []byte("123"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn: conn,
				key:  "kstruct",
				path: ".",
			},
			wantRes: []byte("{\"name\":\"Item#1\",\"number\":1}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONGet(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONGet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONDel(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "kstruct", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				conn: conn,
				key:  "kstr",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldOK",
			args: args{
				conn: conn,
				key:  "kstruct",
				path: "name",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldNotOK",
			args: args{
				conn: conn,
				key:  "kstruct",
				path: "foobar",
			},
			wantRes: int64(0),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn: conn,
				key:  "kstruct",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONDel(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONDel() = %t, want %t", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONMGet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj1 := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	testObj2 := TestObject{
		Name:   "Item#2",
		Number: 2,
	}

	testObj3 := TestObject{
		Name:   "Item#3",
		Number: 3,
	}

	_, err = JSONSet(conn, "testObj1", ".", testObj1, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "testObj2", ".", testObj2, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "testObj3", ".", testObj3, false, false)
	if err != nil {
		return
	}

	resultNameThreeStudents := make([]interface{}, 0)
	type args struct {
		conn redis.Conn
		path string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "NameThreeStudents",
			args: args{
				conn: conn,
				path: "name",
				keys: []string{"testObj1", "testObj2", "testObj3"},
			},
			wantRes: append(resultNameThreeStudents,
				[]byte("\"Item#1\""),
				[]byte("\"Item#2\""),
				[]byte("\"Item#3\""),
			),
			wantErr: false,
		},
		{
			name: "NonExistingKey",
			args: args{
				conn: conn,
				path: "name",
				keys: []string{"testObj1", "testObj2", "foobar"},
			},
			wantRes: append(resultNameThreeStudents,
				[]byte("\"Item#1\""),
				[]byte("\"Item#2\""),
				nil,
			),
			wantErr: false,
		},
		{
			name: "NonExistingKey",
			args: args{
				conn: conn,
				path: "foobar",
				keys: []string{"testObj1"},
			},
			wantRes: append(resultNameThreeStudents,
				nil,
			),
			wantErr: false,
		},
		{
			name: "NoKeys",
			args: args{
				conn: conn,
				path: "name",
				keys: []string{},
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONMGet(tt.args.conn, tt.args.path, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONMGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONMGet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestUnsupportedCommand(t *testing.T) {
	_, _, err := CommandBuilder("FOOBAR", nil)
	if err == nil {
		t.Errorf("TestUnsupportedCommand() returned nil error")
		return
	}
}

func TestJSONType(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "testObj", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "Object",
			args: args{
				conn: conn,
				key:  "testObj",
				path: ".",
			},
			wantRes: "object",
			wantErr: false,
		},
		{
			name: "String",
			args: args{
				conn: conn,
				key:  "testObj",
				path: "name",
			},
			wantRes: "string",
			wantErr: false,
		},
		{
			name: "Integer",
			args: args{
				conn: conn,
				key:  "testObj",
				path: "number",
			},
			wantRes: "integer",
			wantErr: false,
		},
		{
			name: "NotExist",
			args: args{
				conn: conn,
				key:  "foobar",
				path: "number",
			},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONType(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONType() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONNumIncrBy(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kint", ".", 1, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "testObj", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn   redis.Conn
		key    string
		path   string
		number int
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleInt",
			args: args{
				conn:   conn,
				key:    "kint",
				path:   ".",
				number: 5,
			},
			wantRes: []byte("6"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn:   conn,
				key:    "testObj",
				path:   ".number",
				number: 5,
			},
			wantRes: []byte("6"),
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:   conn,
				key:    "kstr",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found string"),
			wantErr: true,
		},
		{
			name: "SimpleStructNotOK",
			args: args{
				conn:   conn,
				key:    "testObj",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found object"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONNumIncrBy(tt.args.conn, tt.args.key, tt.args.path, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONNumIncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONNumIncrBy() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONNumMultBy(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kint", ".", 2, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 2,
	}

	_, err = JSONSet(conn, "testObj", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn   redis.Conn
		key    string
		path   string
		number int
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleInt",
			args: args{
				conn:   conn,
				key:    "kint",
				path:   ".",
				number: 5,
			},
			wantRes: []byte("10"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn:   conn,
				key:    "testObj",
				path:   ".number",
				number: 5,
			},
			wantRes: []byte("10"),
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:   conn,
				key:    "kstr",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found string"),
			wantErr: true,
		},
		{
			name: "SimpleStructNotOK",
			args: args{
				conn:   conn,
				key:    "testObj",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found object"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONNumMultBy(tt.args.conn, tt.args.key, tt.args.path, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONNumMultBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONNumMultBy() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONStrAppend(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "testObj", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn       redis.Conn
		key        string
		path       string
		jsonstring string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				conn:       conn,
				key:        "kstr",
				path:       ".",
				jsonstring: "\"Appended\"",
			},
			wantRes: int64(20),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn:       conn,
				key:        "testObj",
				path:       "name",
				jsonstring: "\"24\"",
			},
			wantRes: int64(8),
			wantErr: false,
		},
		{
			name: "SimpleStructNotOK",
			args: args{
				conn:       conn,
				key:        "testObj",
				path:       "number",
				jsonstring: "\"24\"",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected string but found integer"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONStrAppend(tt.args.conn, tt.args.key, tt.args.path, tt.args.jsonstring)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONStrAppend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONStrAppend() = %v, want %t", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONStrLen(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "testObj", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				conn: conn,
				key:  "kstr",
				path: ".",
			},
			wantRes: int64(12),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn: conn,
				key:  "testObj",
				path: "name",
			},
			wantRes: int64(6),
			wantErr: false,
		},
		{
			name: "SimpleStructNotOK",
			args: args{
				conn: conn,
				key:  "testObj",
				path: "number",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected string but found integer"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONStrLen(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONStrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONStrLen() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONArrAppend(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"one"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}
	appendValues := make([]interface{}, 0)
	valuesAppendStr := []string{"two", "three"}
	for _, valueAppend := range valuesAppendStr {
		appendValues = append(appendValues, valueAppend)
	}

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}

	type args struct {
		conn   redis.Conn
		key    string
		path   string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleArray",
			args: args{
				conn:   conn,
				key:    "karr",
				path:   ".",
				values: appendValues,
			},
			wantRes: int64(3),
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:   conn,
				key:    "kstr",
				path:   ".",
				values: appendValues,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONArrAppend(tt.args.conn, tt.args.key, tt.args.path, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrAppend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONArrAppend() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONArrLen(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}
	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleArray",
			args: args{
				conn: conn,
				key:  "karr",
				path: ".",
			},
			wantRes: int64(3),
			wantErr: false,
		},
		{
			name: "SimpleString",
			args: args{
				conn: conn,
				key:  "kstr",
				path: ".",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONArrLen(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONArrLen() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONArrPop(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three", "four"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}

	type args struct {
		conn   redis.Conn
		key    string
		path   string
		index  int
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleArrayLastPop",
			args: args{
				conn:  conn,
				key:   "karr",
				path:  ".",
				index: PopArrLast,
			},
			wantRes: string("four"),
			wantErr: false,
		},
		{
			name: "SimpleArray2ndElementPop",
			args: args{
				conn:  conn,
				key:   "karr",
				path:  ".",
				index: 1,
			},
			wantRes: "two",
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:  conn,
				key:   "kstr",
				path:  ".",
				index: PopArrLast,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := JSONArrPop(tt.args.conn, tt.args.key, tt.args.path, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var gotRes interface{}
				err = json.Unmarshal(res.([]byte), &gotRes)
				if err != nil {
					t.Errorf("JSONArrPop(): Failed to JSON Unmarshal")
					return
				}
				if !reflect.DeepEqual(gotRes, tt.wantRes) {
					t.Errorf("JSONArrPop() = %v, want %v", gotRes, tt.wantRes)
				}
			}
		})
	}
}

func TestJSONArrIndex(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}

	type args struct {
		conn          redis.Conn
		key           string
		path          string
		value         interface{}
		optionalRange []int
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleArray",
			args: args{
				conn:  conn,
				key:   "karr",
				path:  ".",
				value: "two",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleArrayElementNotPresent",
			args: args{
				conn:  conn,
				key:   "karr",
				path:  ".",
				value: "ten",
			},
			wantRes: int64(-1),
			wantErr: false,
		},
		{
			name: "SimpleArrayElementOutOfRangeWithStart",
			args: args{
				conn:          conn,
				key:           "karr",
				path:          ".",
				value:         "two",
				optionalRange: []int{2},
			},
			wantRes: int64(-1),
			wantErr: false,
		},
		{
			name: "SimpleArrayElementOutOfRange",
			args: args{
				conn:          conn,
				key:           "karr",
				path:          ".",
				value:         "three",
				optionalRange: []int{1, 5},
			},
			wantRes: int64(2),
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:  conn,
				key:   "kstr",
				path:  ".",
				value: "one",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotRes, err := JSONArrIndex(tt.args.conn, tt.args.key, tt.args.path, tt.args.value, tt.args.optionalRange...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotRes.(int64), tt.wantRes.(int64)) {
					t.Errorf("JSONArrIndex() = %v, want %v", gotRes, tt.wantRes)
				}
			}
		})
	}
}

func TestJSONArrTrim(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three", "four"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}

	type args struct {
		conn  redis.Conn
		key   string
		path  string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleArray",
			args: args{
				conn:  conn,
				key:   "karr",
				path:  ".",
				start: 1,
				end:   2,
			},
			wantRes: int64(2),
			wantErr: false,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:  conn,
				key:   "kstr",
				path:  ".",
				start: 1,
				end:   2,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotRes, err := JSONArrTrim(tt.args.conn, tt.args.key, tt.args.path, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrTrim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONArrTrim() = %v, want %v", gotRes, tt.wantRes)
			}

		})
	}
}

func TestJSONArrInsert(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Do("FLUSHALL")
		conn.Close()
	}()

	values := make([]interface{}, 0)
	valuesStr := []string{"three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = JSONSet(conn, "karr", ".", values, false, false)
	if err != nil {
		return
	}
	insertValues := make([]interface{}, 0)
	valuesInsertStr := []string{"one", "two"}
	for _, valueAppend := range valuesInsertStr {
		insertValues = append(insertValues, valueAppend)
	}

	finalStrSlice := make([]string, 0, 4)
	finalStrSlice = append(finalStrSlice, valuesInsertStr...)
	finalStrSlice = append(finalStrSlice, valuesStr...)

	_, err = JSONSet(conn, "kstr", ".", "SimpleString", false, false)
	if err != nil {
		return
	}

	fssM, err := json.Marshal(finalStrSlice)
	if err != nil {
		return
	}

	type args struct {
		conn   redis.Conn
		key    string
		path   string
		index  int
		values []interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantRes       interface{}
		wantErr       bool
		finalSliceGot []byte
	}{
		{
			name: "SimpleArray",
			args: args{
				conn:   conn,
				key:    "karr",
				path:   ".",
				index:  0,
				values: insertValues,
			},
			wantRes:       int64(3),
			wantErr:       false,
			finalSliceGot: fssM,
		},
		{
			name: "SimpleStringNotOK",
			args: args{
				conn:   conn,
				key:    "kstr",
				path:   ".",
				index:  0,
				values: insertValues,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONArrInsert(tt.args.conn, tt.args.key, tt.args.path, tt.args.index, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONArrInsert() = %v, want %v", gotRes, tt.wantRes)
				return
			}

			if !tt.wantErr {
				newArr, err := JSONGet(tt.args.conn, tt.args.key, tt.args.path)
				if err != nil {
					t.Errorf("JSONArrGet(): Failed to JSONGet")
					return
				}
				if !reflect.DeepEqual(newArr.([]byte), tt.finalSliceGot) {
					t.Errorf("JSONArrGet() = %v, want %v", newArr.([]byte), tt.finalSliceGot)
				}
			}
		})
	}
}
