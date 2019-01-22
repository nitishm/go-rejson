package rejson

import (
	"encoding/json"
	"github.com/Shivam010/go-rejson/rjs"
	"reflect"
	"testing"

	"github.com/gomodule/redigo/redis"
)

type TestObject struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func TestUnsupportedCommand(t *testing.T) {
	_, _, err := rjs.CommandBuilder(1234, nil)
	if err == nil {
		t.Errorf("TestUnsupportedCommand() returned nil error")
		return
	}
}

func TestJSONSet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	testObj := TestObject{
		"item#1",
		1,
	}
	type args struct {
		key  string
		path string
		obj  interface{}
		opt  []rjs.SetOption
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
				opt:  []rjs.SetOption{rjs.SetOption_NX},
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
				opt:  []rjs.SetOption{rjs.SetOption_NX},
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
				opt:  []rjs.SetOption{rjs.SetOption_XX},
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
				opt:  []rjs.SetOption{rjs.SetOption_XX},
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
				obj:  "client",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONSet(tt.args.key, tt.args.path, tt.args.obj, tt.args.opt...)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kint", ".", 123)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("kstruct", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
		key     string
		path    string
		options []rjs.GetOption
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
				key:     "kstr",
				path:    ".",
				options: []rjs.GetOption{},
			},
			wantRes: []byte("\"simplestring\""),
			wantErr: false,
		},
		{
			name: "SimpleInt",
			args: args{
				key:     "kint",
				path:    ".",
				options: []rjs.GetOption{},
			},
			wantRes: []byte("123"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:     "kstruct",
				path:    ".",
				options: []rjs.GetOption{},
			},
			wantRes: []byte("{\"name\":\"Item#1\",\"number\":1}"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "kstruct",
				path: ".",
				options: []rjs.GetOption{
					rjs.GETOption_INDENT,
					rjs.GETOption_NEWLINE,
					rjs.GETOption_NOESCAPE,
					rjs.GETOption_SPACE,
				},
			},
			wantRes: []byte("{\n\t\"name\": \"Item#1\",\n\t\"number\": 1\n}"),
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONGet(tt.args.key, tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("\nJSONGet() = %v,\nwant      = %v", gotRes, tt.wantRes)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("kstruct", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:  "kstr",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldOK",
			args: args{
				key:  "kstruct",
				path: "name",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldNotOK",
			args: args{
				key:  "kstruct",
				path: "foobar",
			},
			wantRes: int64(0),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "kstruct",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONDel(tt.args.key, tt.args.path)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
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

	_, err = rh.JSONSet("testObj1", ".", testObj1)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("testObj2", ".", testObj2)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("testObj3", ".", testObj3)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	resultNameThreeStudents := make([]interface{}, 0)
	type args struct {
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
				path: "name",
				keys: []string{},
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				keys: []string{"active"},
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONMGet(tt.args.path, tt.args.keys...)
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

func TestJSONType(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("testObj", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:  "testObj",
				path: ".",
			},
			wantRes: "object",
			wantErr: false,
		},
		{
			name: "String",
			args: args{
				key:  "testObj",
				path: "name",
			},
			wantRes: "string",
			wantErr: false,
		},
		{
			name: "Integer",
			args: args{
				key:  "testObj",
				path: "number",
			},
			wantRes: "integer",
			wantErr: false,
		},
		{
			name: "NotExist",
			args: args{
				key:  "foobar",
				path: "number",
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONType(tt.args.key, tt.args.path)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kint", ".", 1)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("testObj", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:    "testObj",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found object"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONNumIncrBy(tt.args.key, tt.args.path, tt.args.number)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kint", ".", 2)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 2,
	}

	_, err = rh.JSONSet("testObj", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:    "testObj",
				path:   ".",
				number: 5,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected a number but found object"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONNumMultBy(tt.args.key, tt.args.path, tt.args.number)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("testObj", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:        "testObj",
				path:       "number",
				jsonstring: "\"24\"",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected string but found integer"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONStrAppend(tt.args.key, tt.args.path, tt.args.jsonstring)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("testObj", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:  "kstr",
				path: ".",
			},
			wantRes: int64(12),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "testObj",
				path: "name",
			},
			wantRes: int64(6),
			wantErr: false,
		},
		{
			name: "SimpleStructNotOK",
			args: args{
				key:  "testObj",
				path: "number",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected string but found integer"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONStrLen(tt.args.key, tt.args.path)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"one"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	appendValues := make([]interface{}, 0)
	valuesAppendStr := []string{"two", "three"}
	for _, valueAppend := range valuesAppendStr {
		appendValues = append(appendValues, valueAppend)
	}

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:    "kstr",
				path:   ".",
				values: appendValues,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONArrAppend(tt.args.key, tt.args.path, tt.args.values...)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	type args struct {
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
				key:  "karr",
				path: ".",
			},
			wantRes: int64(3),
			wantErr: false,
		},
		{
			name: "SimpleString",
			args: args{
				key:  "kstr",
				path: ".",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONArrLen(tt.args.key, tt.args.path)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three", "four"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
		key   string
		path  string
		index int
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
				key:   "karr",
				path:  ".",
				index: rjs.PopArrLast,
			},
			wantRes: string("four"),
			wantErr: false,
		},
		{
			name: "SimpleArray2ndElementPop",
			args: args{
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
				key:   "kstr",
				path:  ".",
				index: rjs.PopArrLast,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			res, err := rh.JSONArrPop(tt.args.key, tt.args.path, tt.args.index)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:   "kstr",
				path:  ".",
				value: "one",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}

			gotRes, err := rh.JSONArrIndex(tt.args.key, tt.args.path, tt.args.value, tt.args.optionalRange...)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"one", "two", "three", "four"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:   "kstr",
				path:  ".",
				start: 1,
				end:   2,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}

			gotRes, err := rh.JSONArrTrim(tt.args.key, tt.args.path, tt.args.start, tt.args.end)
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
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	values := make([]interface{}, 0)
	valuesStr := []string{"three"}
	for _, value := range valuesStr {
		values = append(values, value)
	}
	_, err = rh.JSONSet("karr", ".", values)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
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

	_, err = rh.JSONSet("kstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	fssM, err := json.Marshal(finalStrSlice)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:    "kstr",
				path:   ".",
				index:  0,
				values: insertValues,
			},
			wantRes: redis.Error("ERR wrong type of path value - expected array but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONArrInsert(tt.args.key, tt.args.path, tt.args.index, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONArrInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONArrInsert() = %v, want %v", gotRes, tt.wantRes)
				return
			}

			if !tt.wantErr {
				newArr, err := rh.JSONGet(tt.args.key, tt.args.path)
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

func TestJSONObjLen(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	type Object struct {
		Name      string `json:"name"`
		LastSeen  int64  `json:"lastSeen"`
		LoggedOut bool   `json:"loggedOut"`
	}
	obj := Object{"Leonard Cohen", 1478476800, true}
	_, err = rh.JSONSet("tobj", ".", obj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("tstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	type args struct {
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
			name: "SimpleObject",
			args: args{
				key:  "tobj",
				path: ".",
			},
			wantRes: int64(3),
			wantErr: false,
		},
		{
			name: "SimpleString",
			args: args{
				key:  "tstr",
				path: ".",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected object but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONObjLen(tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONObjLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONObjLen() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONObjKeys(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	type Object struct {
		Name      string `json:"name"`
		LastSeen  int64  `json:"lastSeen"`
		LoggedOut bool   `json:"loggedOut"`
	}
	obj := Object{"Leonard Cohen", 1478476800, true}
	_, err = rh.JSONSet("tobj", ".", obj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	_, err = rh.JSONSet("tstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	type args struct {
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
			name: "SimpleObject",
			args: args{
				key:  "tobj",
				path: ".",
			},
			wantRes: []string{"name", "lastSeen", "loggedOut"},
			wantErr: false,
		},
		{
			name: "SimpleString",
			args: args{
				key:  "tstr",
				path: ".",
			},
			wantRes: redis.Error("ERR wrong type of path value - expected object but found string"),
			wantErr: true,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONObjKeys(tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONObjKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONObjKeys() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONDebug(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("tstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	type args struct {
		subCommand string
		key        string
		path       string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "Debug Help",
			args: args{
				subCommand: rjs.DebugHelpSubcommand,
				key:        "tstr",
				path:       ".",
			},
			wantRes: rjs.DebugHelpOutput,
			wantErr: false,
		},
		{
			name: "Debug Memory",
			args: args{
				subCommand: rjs.DebugMemorySubcommand,
				key:        "tstr",
				path:       ".",
			},
			wantRes: int64(36),
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONDebug(tt.args.subCommand, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONDebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONDebug() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONForget(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("kstr", ".", "simplestring")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = rh.JSONSet("kstruct", ".", testObj)
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}

	type args struct {
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
				key:  "kstr",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldOK",
			args: args{
				key:  "kstruct",
				path: "name",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: "SimpleStructFieldNotOK",
			args: args{
				key:  "kstruct",
				path: "foobar",
			},
			wantRes: int64(0),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "kstruct",
				path: ".",
			},
			wantRes: int64(1),
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			gotRes, err := rh.JSONForget(tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONForget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONForget() = %t, want %t", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONResp(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		_, err = conn.Do("FLUSHALL")
		err = conn.Close()
		if err != nil {
			t.Fatal("Failed to communicate to redis-server")
		}
	}()
	rh := NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err = rh.JSONSet("tstr", ".", "SimpleString")
	if err != nil {
		t.Fatal("Failed to Set key ", err)
		return
	}
	type args struct {
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RESP",
			args: args{
				key:  "tstr",
				path: ".",
			},
			wantErr: false,
		},
		{
			name: rjs.ClientInactive,
			args: args{
				key:  "active",
				path: ".",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == rjs.ClientInactive {
				rh.SetClientInactive()
			}
			_, err := rh.JSONResp(tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
